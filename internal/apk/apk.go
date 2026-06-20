// Package apk provides core functionality for analyzing Android APK files.
// It includes APK parsing, manifest extraction, and security analysis.
package apk

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// APKInfo contains metadata extracted from an APK file.
type APKInfo struct {
	FilePath     string
	FileSize     int64
	PackageName  string
	VersionName  string
	VersionCode  int
	MinSDK       int
	TargetSDK    int
	Permissions  []string
	Activities   []string
	Services     []string
	Receivers    []string
	Providers    []string
	HasNativeLib bool
	Architectures []string
	IsDebuggable bool
	AllowBackup  bool
}

// ManifestAttribute represents an attribute in the Android manifest.
type ManifestAttribute struct {
	Name  string
	Value string
}

// SecurityIssue represents a potential security issue found in the APK.
type SecurityIssue struct {
	Severity    string // HIGH, MEDIUM, LOW, INFO
	Category    string
	Description string
	Details     string
}

// DangerousPermission represents a dangerous Android permission.
type DangerousPermission struct {
	Name        string
	Description string
	Risk        string
}

// Common dangerous permissions that require runtime approval.
// Based on OWASP MASTG and Android official documentation.
// Reference: https://mas.owasp.org/MASTG/tests/android/MASVS-PLATFORM/MASTG-TEST-0024/
// Reference: https://developer.android.com/reference/android/Manifest.permission
var DangerousPermissions = map[string]DangerousPermission{
	// Calendar permissions
	"android.permission.READ_CALENDAR":  {"READ_CALENDAR", "Read calendar events", "Privacy"},
	"android.permission.WRITE_CALENDAR": {"WRITE_CALENDAR", "Modify calendar events", "Privacy"},

	// Camera permissions
	"android.permission.CAMERA": {"CAMERA", "Access device camera", "Privacy"},

	// Contacts permissions
	"android.permission.READ_CONTACTS":  {"READ_CONTACTS", "Read user contacts", "Privacy"},
	"android.permission.WRITE_CONTACTS": {"WRITE_CONTACTS", "Modify user contacts", "Privacy"},
	"android.permission.GET_ACCOUNTS":   {"GET_ACCOUNTS", "Access device accounts", "Privacy"},

	// Location permissions
	"android.permission.ACCESS_FINE_LOCATION":       {"ACCESS_FINE_LOCATION", "Access precise GPS location", "Privacy/Tracking"},
	"android.permission.ACCESS_COARSE_LOCATION":     {"ACCESS_COARSE_LOCATION", "Access approximate location", "Privacy/Tracking"},
	"android.permission.ACCESS_BACKGROUND_LOCATION": {"ACCESS_BACKGROUND_LOCATION", "Access location in background", "Privacy/Tracking"},

	// Microphone permissions
	"android.permission.RECORD_AUDIO": {"RECORD_AUDIO", "Record audio from microphone", "Privacy"},

	// Phone permissions
	"android.permission.READ_PHONE_STATE":      {"READ_PHONE_STATE", "Read phone state and identity (IMEI)", "Privacy"},
	"android.permission.READ_PHONE_NUMBERS":    {"READ_PHONE_NUMBERS", "Read phone numbers", "Privacy"},
	"android.permission.CALL_PHONE":            {"CALL_PHONE", "Make phone calls without user action", "Financial"},
	"android.permission.ANSWER_PHONE_CALLS":    {"ANSWER_PHONE_CALLS", "Answer incoming calls programmatically", "Privacy"},
	"android.permission.READ_CALL_LOG":         {"READ_CALL_LOG", "Read call history", "Privacy"},
	"android.permission.WRITE_CALL_LOG":        {"WRITE_CALL_LOG", "Modify call history", "Privacy"},
	"android.permission.ADD_VOICEMAIL":         {"ADD_VOICEMAIL", "Add voicemail messages", "Privacy"},
	"android.permission.USE_SIP":               {"USE_SIP", "Use SIP service for calls", "Financial"},
	"android.permission.PROCESS_OUTGOING_CALLS": {"PROCESS_OUTGOING_CALLS", "Intercept outgoing calls", "Privacy"},
	"android.permission.ACCEPT_HANDOVER":       {"ACCEPT_HANDOVER", "Accept call handover from another app", "Privacy"},

	// SMS permissions
	"android.permission.SEND_SMS":         {"SEND_SMS", "Send SMS messages", "Financial"},
	"android.permission.RECEIVE_SMS":      {"RECEIVE_SMS", "Receive SMS messages", "Privacy"},
	"android.permission.READ_SMS":         {"READ_SMS", "Read SMS messages", "Privacy"},
	"android.permission.RECEIVE_WAP_PUSH": {"RECEIVE_WAP_PUSH", "Receive WAP push messages", "Privacy"},
	"android.permission.RECEIVE_MMS":      {"RECEIVE_MMS", "Receive MMS messages", "Privacy"},

	// Storage permissions
	"android.permission.READ_EXTERNAL_STORAGE":  {"READ_EXTERNAL_STORAGE", "Read external storage", "Privacy"},
	"android.permission.WRITE_EXTERNAL_STORAGE": {"WRITE_EXTERNAL_STORAGE", "Write to external storage", "Data"},
	"android.permission.ACCESS_MEDIA_LOCATION":  {"ACCESS_MEDIA_LOCATION", "Access photo/video location metadata", "Privacy"},

	// Body sensors permissions
	"android.permission.BODY_SENSORS":            {"BODY_SENSORS", "Access body sensors (heart rate, etc.)", "Health/Privacy"},
	"android.permission.BODY_SENSORS_BACKGROUND": {"BODY_SENSORS_BACKGROUND", "Access body sensors in background", "Health/Privacy"},

	// Activity recognition permissions
	"android.permission.ACTIVITY_RECOGNITION": {"ACTIVITY_RECOGNITION", "Recognize physical activity", "Privacy"},

	// Nearby devices permissions (Android 12+)
	"android.permission.BLUETOOTH_ADVERTISE": {"BLUETOOTH_ADVERTISE", "Advertise to nearby Bluetooth devices", "Privacy"},
	"android.permission.BLUETOOTH_CONNECT":   {"BLUETOOTH_CONNECT", "Connect to paired Bluetooth devices", "Privacy"},
	"android.permission.BLUETOOTH_SCAN":      {"BLUETOOTH_SCAN", "Scan for nearby Bluetooth devices", "Privacy"},
	"android.permission.NEARBY_WIFI_DEVICES": {"NEARBY_WIFI_DEVICES", "Find nearby WiFi devices", "Privacy"},
	"android.permission.UWB_RANGING":         {"UWB_RANGING", "Determine position with Ultra-Wideband", "Privacy"},

	// Notification permission (Android 13+)
	"android.permission.POST_NOTIFICATIONS": {"POST_NOTIFICATIONS", "Send notifications", "Privacy"},

	// Media permissions (Android 13+)
	"android.permission.READ_MEDIA_AUDIO":  {"READ_MEDIA_AUDIO", "Read audio files", "Privacy"},
	"android.permission.READ_MEDIA_VIDEO":  {"READ_MEDIA_VIDEO", "Read video files", "Privacy"},
	"android.permission.READ_MEDIA_IMAGES": {"READ_MEDIA_IMAGES", "Read image files", "Privacy"},

	// Visual media permission (Android 14+)
	"android.permission.READ_MEDIA_VISUAL_USER_SELECTED": {"READ_MEDIA_VISUAL_USER_SELECTED", "Read selected photos/videos", "Privacy"},
}

// AbusivePermission represents a permission commonly abused by malware.
type AbusivePermission struct {
	Permission  string
	Status      string // normal, dangerous, unknown
	Info        string // short description
	Description string // detailed description
}

// MalwarePermissions are top permissions widely abused by known malware.
// These are high-priority indicators of potentially malicious behavior.
// Based on OWASP MASTG and mobile malware research.
// Reference: https://mas.owasp.org/MASTG/tests/android/MASVS-PLATFORM/MASTG-TEST-0024/
var MalwarePermissions = map[string]AbusivePermission{
	// Network permissions - Required for C2 communication and data exfiltration
	"android.permission.INTERNET": {
		Permission:  "android.permission.INTERNET",
		Status:      "normal",
		Info:        "full network access",
		Description: "Allows the app to create network sockets and use custom network protocols. Required for C2 communication and data exfiltration.",
	},
	"android.permission.ACCESS_NETWORK_STATE": {
		Permission:  "android.permission.ACCESS_NETWORK_STATE",
		Status:      "normal",
		Info:        "view network connections",
		Description: "Allows the app to view network connection state. Malware uses this to check connectivity before data exfiltration.",
	},
	"android.permission.ACCESS_WIFI_STATE": {
		Permission:  "android.permission.ACCESS_WIFI_STATE",
		Status:      "normal",
		Info:        "view Wi-Fi status",
		Description: "Allows the app to view WiFi networking information including BSSID. Used for location tracking and network reconnaissance.",
	},
	"android.permission.CHANGE_WIFI_STATE": {
		Permission:  "android.permission.CHANGE_WIFI_STATE",
		Status:      "normal",
		Info:        "change Wi-Fi state",
		Description: "Allows the app to connect/disconnect WiFi and modify WiFi networks. Can be used for network manipulation attacks.",
	},
	"android.permission.CHANGE_NETWORK_STATE": {
		Permission:  "android.permission.CHANGE_NETWORK_STATE",
		Status:      "normal",
		Info:        "change network connectivity",
		Description: "Allows the app to change network connectivity state. Can be used to force traffic through malicious networks.",
	},

	// Persistence permissions - Ensures malware survives reboot
	"android.permission.RECEIVE_BOOT_COMPLETED": {
		Permission:  "android.permission.RECEIVE_BOOT_COMPLETED",
		Status:      "normal",
		Info:        "automatically start at boot",
		Description: "Allows the app to start itself after device boot. Critical for malware persistence across device restarts.",
	},
	"android.permission.WAKE_LOCK": {
		Permission:  "android.permission.WAKE_LOCK",
		Status:      "normal",
		Info:        "prevent phone from sleeping",
		Description: "Allows the app to prevent the device from sleeping. Malware uses this to run continuously in the background.",
	},
	"android.permission.REQUEST_IGNORE_BATTERY_OPTIMIZATIONS": {
		Permission:  "android.permission.REQUEST_IGNORE_BATTERY_OPTIMIZATIONS",
		Status:      "normal",
		Info:        "ignore battery optimizations",
		Description: "Allows the app to request exemption from battery optimizations. Used by malware to maintain persistent background execution.",
	},
	"android.permission.FOREGROUND_SERVICE": {
		Permission:  "android.permission.FOREGROUND_SERVICE",
		Status:      "normal",
		Info:        "run foreground service",
		Description: "Allows the app to use Service.startForeground. Used by malware to maintain persistent background execution with elevated priority.",
	},
	"android.permission.FOREGROUND_SERVICE_DATA_SYNC": {
		Permission:  "android.permission.FOREGROUND_SERVICE_DATA_SYNC",
		Status:      "normal",
		Info:        "foreground data sync service",
		Description: "Allows foreground service with dataSync type. Can be misused for persistent data exfiltration.",
	},

	// SMS permissions - Used for OTP interception and premium fraud
	"android.permission.SEND_SMS": {
		Permission:  "android.permission.SEND_SMS",
		Status:      "dangerous",
		Info:        "send SMS messages",
		Description: "Allows the app to send SMS messages. Used for premium SMS fraud, spreading malware, and financial theft.",
	},
	"android.permission.RECEIVE_SMS": {
		Permission:  "android.permission.RECEIVE_SMS",
		Status:      "dangerous",
		Info:        "receive SMS messages",
		Description: "Allows the app to receive SMS messages. Critical for intercepting OTP codes and 2FA messages for banking trojans.",
	},
	"android.permission.READ_SMS": {
		Permission:  "android.permission.READ_SMS",
		Status:      "dangerous",
		Info:        "read SMS messages",
		Description: "Allows the app to read SMS messages. Used to steal OTP codes, banking information, and personal communications.",
	},
	"android.permission.RECEIVE_MMS": {
		Permission:  "android.permission.RECEIVE_MMS",
		Status:      "dangerous",
		Info:        "receive MMS messages",
		Description: "Allows the app to receive MMS messages. Can be used to intercept multimedia messages containing sensitive data.",
	},
	"android.permission.RECEIVE_WAP_PUSH": {
		Permission:  "android.permission.RECEIVE_WAP_PUSH",
		Status:      "dangerous",
		Info:        "receive WAP push messages",
		Description: "Allows the app to receive WAP push messages. Can be used to intercept service messages.",
	},

	// Phone/Call permissions - Used for surveillance and fraud
	"android.permission.READ_PHONE_STATE": {
		Permission:  "android.permission.READ_PHONE_STATE",
		Status:      "dangerous",
		Info:        "read phone state and identity",
		Description: "Allows access to phone state including IMEI, phone number, and call state. Used for device fingerprinting and tracking.",
	},
	"android.permission.READ_PHONE_NUMBERS": {
		Permission:  "android.permission.READ_PHONE_NUMBERS",
		Status:      "dangerous",
		Info:        "read phone numbers",
		Description: "Allows the app to read device phone numbers. Used for identity theft and targeted attacks.",
	},
	"android.permission.CALL_PHONE": {
		Permission:  "android.permission.CALL_PHONE",
		Status:      "dangerous",
		Info:        "directly call phone numbers",
		Description: "Allows the app to make phone calls without user interaction. Used for premium number fraud and financial theft.",
	},
	"android.permission.ANSWER_PHONE_CALLS": {
		Permission:  "android.permission.ANSWER_PHONE_CALLS",
		Status:      "dangerous",
		Info:        "answer phone calls",
		Description: "Allows the app to answer incoming calls programmatically. Can be used for call interception.",
	},
	"android.permission.READ_CALL_LOG": {
		Permission:  "android.permission.READ_CALL_LOG",
		Status:      "dangerous",
		Info:        "read call log",
		Description: "Allows the app to read call history. Used to profile contacts and identify high-value targets.",
	},
	"android.permission.WRITE_CALL_LOG": {
		Permission:  "android.permission.WRITE_CALL_LOG",
		Status:      "dangerous",
		Info:        "modify call log",
		Description: "Allows the app to modify call history. Can be used to hide malicious call activity.",
	},
	"android.permission.PROCESS_OUTGOING_CALLS": {
		Permission:  "android.permission.PROCESS_OUTGOING_CALLS",
		Status:      "dangerous",
		Info:        "reroute outgoing calls",
		Description: "Allows the app to intercept outgoing calls. Used to redirect calls, enable fraud, or block emergency services.",
	},

	// Location permissions - Used for stalking and surveillance
	"android.permission.ACCESS_FINE_LOCATION": {
		Permission:  "android.permission.ACCESS_FINE_LOCATION",
		Status:      "dangerous",
		Info:        "access precise GPS location",
		Description: "Allows access to precise GPS location. Used for stalking, surveillance, and targeted physical attacks.",
	},
	"android.permission.ACCESS_COARSE_LOCATION": {
		Permission:  "android.permission.ACCESS_COARSE_LOCATION",
		Status:      "dangerous",
		Info:        "access approximate location",
		Description: "Allows access to approximate location via network. Used for tracking user movements and location profiling.",
	},
	"android.permission.ACCESS_BACKGROUND_LOCATION": {
		Permission:  "android.permission.ACCESS_BACKGROUND_LOCATION",
		Status:      "dangerous",
		Info:        "access location in background",
		Description: "Allows continuous background location access. Enables persistent stalking and surveillance without user awareness.",
	},

	// Camera/Audio permissions - Used for surveillance
	"android.permission.CAMERA": {
		Permission:  "android.permission.CAMERA",
		Status:      "dangerous",
		Info:        "take pictures and videos",
		Description: "Allows the app to use the device camera. Malware uses this for covert photo/video surveillance.",
	},
	"android.permission.RECORD_AUDIO": {
		Permission:  "android.permission.RECORD_AUDIO",
		Status:      "dangerous",
		Info:        "record audio",
		Description: "Allows the app to record audio using the microphone. Commonly abused for eavesdropping and audio surveillance.",
	},

	// Contact permissions - Used for social engineering
	"android.permission.READ_CONTACTS": {
		Permission:  "android.permission.READ_CONTACTS",
		Status:      "dangerous",
		Info:        "read contacts",
		Description: "Allows the app to read contact data. Used to harvest contact lists for spam, phishing, and social engineering attacks.",
	},
	"android.permission.WRITE_CONTACTS": {
		Permission:  "android.permission.WRITE_CONTACTS",
		Status:      "dangerous",
		Info:        "modify contacts",
		Description: "Allows the app to modify contacts. Can be used to inject malicious contact information.",
	},
	"android.permission.GET_ACCOUNTS": {
		Permission:  "android.permission.GET_ACCOUNTS",
		Status:      "dangerous",
		Info:        "find accounts on device",
		Description: "Allows the app to get the list of accounts on the device. Used to identify targets and associated financial services.",
	},

	// Storage permissions - Used for data theft
	"android.permission.READ_EXTERNAL_STORAGE": {
		Permission:  "android.permission.READ_EXTERNAL_STORAGE",
		Status:      "dangerous",
		Info:        "read storage contents",
		Description: "Allows the app to read from external storage. Used by malware to steal photos, documents, and sensitive files.",
	},
	"android.permission.WRITE_EXTERNAL_STORAGE": {
		Permission:  "android.permission.WRITE_EXTERNAL_STORAGE",
		Status:      "dangerous",
		Info:        "modify or delete storage contents",
		Description: "Allows the app to write to external storage. Malware uses this to store stolen data, inject files, or deploy ransomware.",
	},
	"android.permission.MANAGE_EXTERNAL_STORAGE": {
		Permission:  "android.permission.MANAGE_EXTERNAL_STORAGE",
		Status:      "special",
		Info:        "manage all files",
		Description: "Allows broad access to all files on device storage. Extremely high-risk permission for data theft and ransomware.",
	},

	// Overlay and system permissions - Used for phishing and credential theft
	"android.permission.SYSTEM_ALERT_WINDOW": {
		Permission:  "android.permission.SYSTEM_ALERT_WINDOW",
		Status:      "special",
		Info:        "draw over other apps",
		Description: "Allows the app to draw overlays on top of other apps. Critical for overlay attacks, phishing, and credential theft in banking trojans.",
	},
	"android.permission.WRITE_SETTINGS": {
		Permission:  "android.permission.WRITE_SETTINGS",
		Status:      "special",
		Info:        "modify system settings",
		Description: "Allows the app to modify system settings. Can be used to weaken security settings or enable malicious configurations.",
	},

	// Package management - Used to install additional malware
	"android.permission.REQUEST_INSTALL_PACKAGES": {
		Permission:  "android.permission.REQUEST_INSTALL_PACKAGES",
		Status:      "special",
		Info:        "request install packages",
		Description: "Allows the app to request installation of packages. Used by droppers to install additional malware components.",
	},
	"android.permission.REQUEST_DELETE_PACKAGES": {
		Permission:  "android.permission.REQUEST_DELETE_PACKAGES",
		Status:      "normal",
		Info:        "request delete packages",
		Description: "Allows the app to request deletion of packages. Can be used to remove security apps or competing malware.",
	},
	"android.permission.QUERY_ALL_PACKAGES": {
		Permission:  "android.permission.QUERY_ALL_PACKAGES",
		Status:      "normal",
		Info:        "query all packages",
		Description: "Allows the app to see all installed packages. Used for reconnaissance to identify banking apps and security software.",
	},

	// Accessibility and admin - Most dangerous for on-device fraud
	"android.permission.BIND_ACCESSIBILITY_SERVICE": {
		Permission:  "android.permission.BIND_ACCESSIBILITY_SERVICE",
		Status:      "signature",
		Info:        "bind to accessibility service",
		Description: "Allows binding to accessibility services. CRITICAL: Enables full device control, keylogging, screen capture, and automated actions for banking trojans.",
	},
	"android.permission.BIND_DEVICE_ADMIN": {
		Permission:  "android.permission.BIND_DEVICE_ADMIN",
		Status:      "signature",
		Info:        "interact with device admin",
		Description: "Allows interaction with device administrator. Used by ransomware to lock devices and resist uninstallation.",
	},
	"android.permission.BIND_NOTIFICATION_LISTENER_SERVICE": {
		Permission:  "android.permission.BIND_NOTIFICATION_LISTENER_SERVICE",
		Status:      "signature",
		Info:        "bind to notification listener",
		Description: "Allows the app to read all notifications. Critical for intercepting OTPs, banking alerts, and 2FA codes.",
	},

	// Other high-risk permissions
	"android.permission.VIBRATE": {
		Permission:  "android.permission.VIBRATE",
		Status:      "normal",
		Info:        "control vibration",
		Description: "Allows the app to control the vibrator. Can be used to silently acknowledge commands from C2 servers.",
	},
	"android.permission.DISABLE_KEYGUARD": {
		Permission:  "android.permission.DISABLE_KEYGUARD",
		Status:      "normal",
		Info:        "disable keyguard",
		Description: "Allows the app to disable the keyguard (lock screen). Can be exploited to bypass device security.",
	},
	"android.permission.READ_CALENDAR": {
		Permission:  "android.permission.READ_CALENDAR",
		Status:      "dangerous",
		Info:        "read calendar events",
		Description: "Allows reading calendar events. Used to gather intelligence about user's schedule and meetings.",
	},
	"android.permission.WRITE_CALENDAR": {
		Permission:  "android.permission.WRITE_CALENDAR",
		Status:      "dangerous",
		Info:        "modify calendar events",
		Description: "Allows modifying calendar events. Can be used to inject fake appointments for social engineering.",
	},
	"android.permission.USE_BIOMETRIC": {
		Permission:  "android.permission.USE_BIOMETRIC",
		Status:      "normal",
		Info:        "use biometric hardware",
		Description: "Allows use of biometric hardware. Can be misused to bypass or manipulate biometric authentication.",
	},
	"android.permission.USE_FINGERPRINT": {
		Permission:  "android.permission.USE_FINGERPRINT",
		Status:      "normal",
		Info:        "use fingerprint hardware",
		Description: "Deprecated in API 28. Allows use of fingerprint hardware for authentication.",
	},
}

// OtherCommonAbusedPermissions are permissions commonly abused by known malware,
// but are also frequently found in legitimate apps. Requires contextual analysis.
// Based on OWASP MASTG and mobile security research.
var OtherCommonAbusedPermissions = map[string]AbusivePermission{
	// Audio/Media permissions
	"android.permission.MODIFY_AUDIO_SETTINGS": {
		Permission:  "android.permission.MODIFY_AUDIO_SETTINGS",
		Status:      "normal",
		Info:        "change audio settings",
		Description: "Allows the app to modify global audio settings. Can be used to mute the device during malicious activity.",
	},

	// Foreground service permissions - Can enable persistent background execution
	"android.permission.FOREGROUND_SERVICE_CAMERA": {
		Permission:  "android.permission.FOREGROUND_SERVICE_CAMERA",
		Status:      "normal",
		Info:        "foreground service with camera",
		Description: "Allows foreground service with camera type. Can enable covert camera access while appearing as a legitimate service.",
	},
	"android.permission.FOREGROUND_SERVICE_MICROPHONE": {
		Permission:  "android.permission.FOREGROUND_SERVICE_MICROPHONE",
		Status:      "normal",
		Info:        "foreground service with microphone",
		Description: "Allows foreground service with microphone type. Can enable covert audio recording while appearing as a legitimate service.",
	},
	"android.permission.FOREGROUND_SERVICE_MEDIA_PLAYBACK": {
		Permission:  "android.permission.FOREGROUND_SERVICE_MEDIA_PLAYBACK",
		Status:      "normal",
		Info:        "foreground service for media",
		Description: "Allows foreground service with mediaPlayback type. Can be misused for persistent background execution.",
	},
	"android.permission.FOREGROUND_SERVICE_LOCATION": {
		Permission:  "android.permission.FOREGROUND_SERVICE_LOCATION",
		Status:      "normal",
		Info:        "foreground service with location",
		Description: "Allows foreground service with location type. Can enable persistent location tracking.",
	},
	"android.permission.FOREGROUND_SERVICE_PHONE_CALL": {
		Permission:  "android.permission.FOREGROUND_SERVICE_PHONE_CALL",
		Status:      "normal",
		Info:        "foreground service for calls",
		Description: "Allows foreground service with phoneCall type. Can be used for call monitoring.",
	},
	"android.permission.FOREGROUND_SERVICE_CONNECTED_DEVICE": {
		Permission:  "android.permission.FOREGROUND_SERVICE_CONNECTED_DEVICE",
		Status:      "normal",
		Info:        "foreground service for devices",
		Description: "Allows foreground service with connectedDevice type. Can be used for persistent device monitoring.",
	},
	"android.permission.FOREGROUND_SERVICE_SPECIAL_USE": {
		Permission:  "android.permission.FOREGROUND_SERVICE_SPECIAL_USE",
		Status:      "normal",
		Info:        "foreground service special use",
		Description: "Allows foreground service with specialUse type. Can be abused for arbitrary persistent background execution.",
	},

	// Advertising and tracking permissions
	"com.google.android.c2dm.permission.RECEIVE": {
		Permission:  "com.google.android.c2dm.permission.RECEIVE",
		Status:      "normal",
		Info:        "receive push notifications",
		Description: "Allows receiving push notifications from cloud. Can be used as a C2 channel for receiving malware commands.",
	},
	"com.google.android.finsky.permission.BIND_GET_INSTALL_REFERRER_SERVICE": {
		Permission:  "com.google.android.finsky.permission.BIND_GET_INSTALL_REFERRER_SERVICE",
		Status:      "normal",
		Info:        "bind install referrer service",
		Description: "Google permission for tracking app installation referrers. Used for attribution tracking.",
	},
	"com.google.android.gms.permission.AD_ID": {
		Permission:  "com.google.android.gms.permission.AD_ID",
		Status:      "normal",
		Info:        "access advertising ID",
		Description: "Allows access to Google advertising ID. Can be used for cross-app user tracking.",
	},
	"android.permission.ACCESS_ADSERVICES_ATTRIBUTION": {
		Permission:  "android.permission.ACCESS_ADSERVICES_ATTRIBUTION",
		Status:      "normal",
		Info:        "access ad attribution",
		Description: "Allows access to advertising attribution data. Can gather data about user ad interactions.",
	},
	"android.permission.ACCESS_ADSERVICES_AD_ID": {
		Permission:  "android.permission.ACCESS_ADSERVICES_AD_ID",
		Status:      "normal",
		Info:        "access advertising ID",
		Description: "Allows access to device advertising ID for user tracking across apps.",
	},
	"android.permission.ACCESS_ADSERVICES_TOPICS": {
		Permission:  "android.permission.ACCESS_ADSERVICES_TOPICS",
		Status:      "normal",
		Info:        "access ad topics",
		Description: "Allows access to advertising topics/interests for targeted advertising and profiling.",
	},

	// Download and file access
	"android.permission.DOWNLOAD_WITHOUT_NOTIFICATION": {
		Permission:  "android.permission.DOWNLOAD_WITHOUT_NOTIFICATION",
		Status:      "normal",
		Info:        "silent downloads",
		Description: "Allows downloading files without showing notification. Used for covert malware downloads.",
	},
	"android.permission.ACCESS_MEDIA_LOCATION": {
		Permission:  "android.permission.ACCESS_MEDIA_LOCATION",
		Status:      "dangerous",
		Info:        "access photo/video location",
		Description: "Allows access to location metadata in photos/videos. Can reveal user's location history.",
	},

	// Bluetooth permissions (Android 12+)
	"android.permission.BLUETOOTH": {
		Permission:  "android.permission.BLUETOOTH",
		Status:      "normal",
		Info:        "access Bluetooth",
		Description: "Allows pairing with Bluetooth devices. Can be used for proximity tracking.",
	},
	"android.permission.BLUETOOTH_ADMIN": {
		Permission:  "android.permission.BLUETOOTH_ADMIN",
		Status:      "normal",
		Info:        "Bluetooth admin",
		Description: "Allows discovering and pairing Bluetooth devices. Can enable device tracking.",
	},
	"android.permission.BLUETOOTH_CONNECT": {
		Permission:  "android.permission.BLUETOOTH_CONNECT",
		Status:      "dangerous",
		Info:        "connect to Bluetooth",
		Description: "Allows connecting to paired Bluetooth devices. Can access sensitive paired device information.",
	},
	"android.permission.BLUETOOTH_SCAN": {
		Permission:  "android.permission.BLUETOOTH_SCAN",
		Status:      "dangerous",
		Info:        "scan for Bluetooth",
		Description: "Allows scanning for nearby Bluetooth devices. Can be used for location tracking.",
	},

	// NFC permission
	"android.permission.NFC": {
		Permission:  "android.permission.NFC",
		Status:      "normal",
		Info:        "use NFC",
		Description: "Allows NFC communication. Could potentially be used for payment-related attacks.",
	},

	// Body sensors
	"android.permission.BODY_SENSORS": {
		Permission:  "android.permission.BODY_SENSORS",
		Status:      "dangerous",
		Info:        "access body sensors",
		Description: "Allows access to body sensors like heart rate monitors. Can collect sensitive health data.",
	},
	"android.permission.ACTIVITY_RECOGNITION": {
		Permission:  "android.permission.ACTIVITY_RECOGNITION",
		Status:      "dangerous",
		Info:        "recognize activity",
		Description: "Allows recognition of physical activity. Can track user movements and behavior patterns.",
	},

	// Usage stats and notification access
	"android.permission.PACKAGE_USAGE_STATS": {
		Permission:  "android.permission.PACKAGE_USAGE_STATS",
		Status:      "special",
		Info:        "access usage stats",
		Description: "Allows access to app usage statistics. Can profile user behavior and app usage patterns.",
	},
	"android.permission.ACCESS_NOTIFICATION_POLICY": {
		Permission:  "android.permission.ACCESS_NOTIFICATION_POLICY",
		Status:      "normal",
		Info:        "access DND settings",
		Description: "Allows access to Do Not Disturb notification policy. Can suppress security notifications.",
	},

	// Alarm and scheduling
	"android.permission.SCHEDULE_EXACT_ALARM": {
		Permission:  "android.permission.SCHEDULE_EXACT_ALARM",
		Status:      "normal",
		Info:        "schedule exact alarms",
		Description: "Allows scheduling exact alarms. Can be used for timed malicious activities.",
	},
	"android.permission.USE_EXACT_ALARM": {
		Permission:  "android.permission.USE_EXACT_ALARM",
		Status:      "normal",
		Info:        "use exact alarms",
		Description: "Allows using exact alarms for time-sensitive features. Can schedule precise malicious actions.",
	},

	// VPN permission
	"android.permission.BIND_VPN_SERVICE": {
		Permission:  "android.permission.BIND_VPN_SERVICE",
		Status:      "signature",
		Info:        "bind VPN service",
		Description: "Allows binding to VPN service. Can intercept all network traffic if misused.",
	},

	// Input method permission
	"android.permission.BIND_INPUT_METHOD": {
		Permission:  "android.permission.BIND_INPUT_METHOD",
		Status:      "signature",
		Info:        "bind input method",
		Description: "Allows binding to input method service. Can enable keylogging if misused.",
	},

	// Clipboard access (Android 10+)
	"android.permission.READ_CLIPBOARD": {
		Permission:  "android.permission.READ_CLIPBOARD",
		Status:      "normal",
		Info:        "read clipboard",
		Description: "Allows reading clipboard content. Can capture passwords, OTPs, and sensitive data copied by user.",
	},

	// Telephony permissions
	"android.permission.READ_PRECISE_PHONE_STATE": {
		Permission:  "android.permission.READ_PRECISE_PHONE_STATE",
		Status:      "signature",
		Info:        "read precise phone state",
		Description: "Allows reading precise phone state information. Provides detailed call state for surveillance.",
	},

	// Flashlight/Camera hardware
	"android.permission.FLASHLIGHT": {
		Permission:  "android.permission.FLASHLIGHT",
		Status:      "normal",
		Info:        "control flashlight",
		Description: "Allows control of camera flashlight. Can be used to signal or draw attention.",
	},

	// System features
	"android.permission.EXPAND_STATUS_BAR": {
		Permission:  "android.permission.EXPAND_STATUS_BAR",
		Status:      "normal",
		Info:        "expand status bar",
		Description: "Allows expanding/collapsing status bar. Can be used to manipulate UI for phishing.",
	},
	"android.permission.GET_TASKS": {
		Permission:  "android.permission.GET_TASKS",
		Status:      "normal",
		Info:        "get running tasks",
		Description: "Deprecated but still used. Allows getting information about running tasks for app profiling.",
	},
	"android.permission.REORDER_TASKS": {
		Permission:  "android.permission.REORDER_TASKS",
		Status:      "normal",
		Info:        "reorder tasks",
		Description: "Allows reordering tasks to foreground. Can be used to manipulate UI state.",
	},
	"android.permission.KILL_BACKGROUND_PROCESSES": {
		Permission:  "android.permission.KILL_BACKGROUND_PROCESSES",
		Status:      "normal",
		Info:        "kill background processes",
		Description: "Allows killing background processes. Can be used to terminate security apps.",
	},

	// Add voicemail
	"android.permission.ADD_VOICEMAIL": {
		Permission:  "android.permission.ADD_VOICEMAIL",
		Status:      "dangerous",
		Info:        "add voicemail",
		Description: "Allows adding voicemail messages. Can be used for voice phishing attacks.",
	},

	// SIP calling
	"android.permission.USE_SIP": {
		Permission:  "android.permission.USE_SIP",
		Status:      "dangerous",
		Info:        "use SIP service",
		Description: "Allows using SIP service for calls. Can be used for call fraud.",
	},

	// Media permissions (Android 13+)
	"android.permission.READ_MEDIA_AUDIO": {
		Permission:  "android.permission.READ_MEDIA_AUDIO",
		Status:      "dangerous",
		Info:        "read audio files",
		Description: "Allows reading audio files. Can access voice recordings and sensitive audio.",
	},
	"android.permission.READ_MEDIA_VIDEO": {
		Permission:  "android.permission.READ_MEDIA_VIDEO",
		Status:      "dangerous",
		Info:        "read video files",
		Description: "Allows reading video files. Can access private video recordings.",
	},
	"android.permission.READ_MEDIA_IMAGES": {
		Permission:  "android.permission.READ_MEDIA_IMAGES",
		Status:      "dangerous",
		Info:        "read image files",
		Description: "Allows reading image files. Can access private photos and screenshots.",
	},

	// Notification permission (Android 13+)
	"android.permission.POST_NOTIFICATIONS": {
		Permission:  "android.permission.POST_NOTIFICATIONS",
		Status:      "dangerous",
		Info:        "send notifications",
		Description: "Allows posting notifications. Can be used for persistent user interaction and social engineering.",
	},
}

// AbusivePermissionResult represents the analysis result for abusive permissions.
type AbusivePermissionResult struct {
	TotalPermissions int
	MalwareMatches   []AbusivePermission
	OtherMatches     []AbusivePermission
	AllPermissions   []string
}

// AnalyzeAbusivePermissions checks APK permissions against known abusive permission lists.
func AnalyzeAbusivePermissions(apkPath string) (*AbusivePermissionResult, error) {
	info, err := ParseManifestBasic(apkPath)
	if err != nil {
		return nil, err
	}

	result := &AbusivePermissionResult{
		TotalPermissions: len(info.Permissions),
		AllPermissions:   info.Permissions,
	}

	for _, perm := range info.Permissions {
		if abusive, found := MalwarePermissions[perm]; found {
			result.MalwareMatches = append(result.MalwareMatches, abusive)
		}
		if abusive, found := OtherCommonAbusedPermissions[perm]; found {
			result.OtherMatches = append(result.OtherMatches, abusive)
		}
	}

	return result, nil
}

// OpenAPK opens an APK file and returns a zip reader.
func OpenAPK(path string) (*zip.ReadCloser, error) {
	return zip.OpenReader(path)
}

// GetAPKInfo extracts basic information from an APK file.
func GetAPKInfo(path string) (*APKInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %v", err)
	}

	apkInfo := &APKInfo{
		FilePath: path,
		FileSize: info.Size(),
	}

	zr, err := OpenAPK(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open APK: %v", err)
	}
	defer zr.Close()

	// Check for native libraries
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "lib/") && strings.HasSuffix(f.Name, ".so") {
			apkInfo.HasNativeLib = true
			// Extract architecture
			parts := strings.Split(f.Name, "/")
			if len(parts) >= 2 {
				arch := parts[1]
				found := false
				for _, a := range apkInfo.Architectures {
					if a == arch {
						found = true
						break
					}
				}
				if !found {
					apkInfo.Architectures = append(apkInfo.Architectures, arch)
				}
			}
		}
	}

	return apkInfo, nil
}

// ListFiles returns all files in the APK.
func ListFiles(path string) ([]string, error) {
	zr, err := OpenAPK(path)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	var files []string
	for _, f := range zr.File {
		files = append(files, f.Name)
	}
	return files, nil
}

// ExtractFile extracts a specific file from the APK.
func ExtractFile(apkPath, fileName, outputPath string) error {
	zr, err := OpenAPK(apkPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == fileName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			// Create output directory if needed
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			return err
		}
	}

	return fmt.Errorf("file not found in APK: %s", fileName)
}

// ExtractStrings extracts readable strings from a file in the APK.
func ExtractStrings(apkPath, fileName string, minLength int) ([]string, error) {
	zr, err := OpenAPK(apkPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == fileName {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			return extractStringsFromBytes(data, minLength), nil
		}
	}

	return nil, fmt.Errorf("file not found: %s", fileName)
}

// ExtractAllStrings extracts all readable strings from the APK.
func ExtractAllStrings(apkPath string, minLength int) (map[string][]string, error) {
	zr, err := OpenAPK(apkPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	results := make(map[string][]string)

	for _, f := range zr.File {
		// Skip directories and large files
		if f.FileInfo().IsDir() || f.UncompressedSize64 > 10*1024*1024 {
			continue
		}

		// Focus on interesting files
		if strings.HasSuffix(f.Name, ".dex") ||
			strings.HasSuffix(f.Name, ".so") ||
			strings.HasSuffix(f.Name, ".xml") ||
			strings.HasSuffix(f.Name, ".json") ||
			strings.HasSuffix(f.Name, ".txt") ||
			strings.HasSuffix(f.Name, ".properties") {

			rc, err := f.Open()
			if err != nil {
				continue
			}

			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			strings := extractStringsFromBytes(data, minLength)
			if len(strings) > 0 {
				results[f.Name] = strings
			}
		}
	}

	return results, nil
}

// extractStringsFromBytes extracts printable strings from binary data.
func extractStringsFromBytes(data []byte, minLength int) []string {
	var strings []string
	var current []byte

	for _, b := range data {
		if b >= 32 && b < 127 {
			current = append(current, b)
		} else {
			if len(current) >= minLength {
				strings = append(strings, string(current))
			}
			current = nil
		}
	}

	if len(current) >= minLength {
		strings = append(strings, string(current))
	}

	return strings
}

// SearchStrings searches for strings matching a pattern in the APK.
func SearchStrings(apkPath, pattern string) (map[string][]string, error) {
	allStrings, err := ExtractAllStrings(apkPath, 4)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	results := make(map[string][]string)
	for file, strs := range allStrings {
		var matches []string
		for _, s := range strs {
			if re.MatchString(s) {
				matches = append(matches, s)
			}
		}
		if len(matches) > 0 {
			results[file] = matches
		}
	}

	return results, nil
}

// AnalyzeSecurity performs security analysis on the APK.
// Based on OWASP MASTG security testing guidelines.
// Reference: https://mas.owasp.org/MASTG/
func AnalyzeSecurity(apkPath string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	info, err := GetAPKInfo(apkPath)
	if err != nil {
		return nil, err
	}

	// Check for debuggable flag (would need manifest parsing)
	if info.IsDebuggable {
		issues = append(issues, SecurityIssue{
			Severity:    "HIGH",
			Category:    "Configuration",
			Description: "Application is debuggable",
			Details:     "The android:debuggable flag is set to true. This allows attackers to attach debuggers and inspect app internals.",
		})
	}

	// Check for backup flag
	if info.AllowBackup {
		issues = append(issues, SecurityIssue{
			Severity:    "MEDIUM",
			Category:    "Configuration",
			Description: "Application allows backup",
			Details:     "The android:allowBackup flag is true. App data can be extracted via ADB backup.",
		})
	}

	// Search for sensitive patterns - Based on OWASP MASTG
	sensitivePatterns := map[string]SecurityIssue{
		// Hardcoded Secrets (MASTG-TEST-0001)
		`(?i)api[_-]?key\s*[=:]`:               {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential API key assignment found"},
		`(?i)api[_-]?secret`:                   {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential API secret found"},
		`(?i)secret[_-]?key`:                   {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential secret key found"},
		`(?i)private[_-]?key`:                  {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential private key found"},
		`(?i)client[_-]?secret`:                {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential client secret found"},
		`(?i)auth[_-]?token`:                   {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential auth token found"},
		`(?i)access[_-]?token`:                 {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential access token found"},
		`(?i)bearer\s+[a-zA-Z0-9\-_]+`:         {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential Bearer token found"},
		`(?i)password\s*[=:]`:                  {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential hardcoded password found"},
		`(?i)passwd\s*[=:]`:                    {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential hardcoded password found"},
		`(?i)encryption[_-]?key`:               {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential encryption key found"},
		`(?i)signing[_-]?key`:                  {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential signing key found"},

		// Cloud Provider Credentials
		`AKIA[0-9A-Z]{16}`:                     {Severity: "CRITICAL", Category: "Cloud Credentials", Description: "AWS Access Key ID found"},
		`(?i)aws[_-]?secret`:                   {Severity: "CRITICAL", Category: "Cloud Credentials", Description: "Potential AWS secret key found"},
		`(?i)aws[_-]?access`:                   {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential AWS access key found"},
		`(?i)aws[_-]?session`:                  {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential AWS session token found"},
		`AIza[0-9A-Za-z\-_]{35}`:               {Severity: "HIGH", Category: "Cloud Credentials", Description: "Google API key found"},
		`(?i)gcp[_-]?api[_-]?key`:              {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential GCP API key found"},
		`(?i)azure[_-]?`:                       {Severity: "MEDIUM", Category: "Cloud Credentials", Description: "Azure reference found"},
		`(?i)digitalocean`:                     {Severity: "MEDIUM", Category: "Cloud Credentials", Description: "DigitalOcean reference found"},
		`(?i)heroku`:                           {Severity: "MEDIUM", Category: "Cloud Credentials", Description: "Heroku reference found"},

		// Firebase Configuration
		`(?i)firebase[_-]?api`:                 {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase API reference found"},
		`(?i)firebase[_-]?url`:                 {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase URL found"},
		`\.firebaseio\.com`:                    {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase Realtime Database URL found"},
		`\.firebaseapp\.com`:                   {Severity: "INFO", Category: "Firebase", Description: "Firebase hosting URL found"},
		`(?i)firebase[_-]?database`:            {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase database reference found"},

		// Third-Party Services
		`(?i)stripe[_-]?`:                      {Severity: "MEDIUM", Category: "Payment", Description: "Stripe payment integration detected"},
		`sk_live_[0-9a-zA-Z]{24}`:              {Severity: "CRITICAL", Category: "Payment", Description: "Stripe live secret key found"},
		`pk_live_[0-9a-zA-Z]{24}`:              {Severity: "HIGH", Category: "Payment", Description: "Stripe live publishable key found"},
		`(?i)paypal`:                           {Severity: "INFO", Category: "Payment", Description: "PayPal integration detected"},
		`(?i)braintree`:                        {Severity: "INFO", Category: "Payment", Description: "Braintree integration detected"},
		`(?i)twilio`:                           {Severity: "MEDIUM", Category: "Third Party", Description: "Twilio integration detected"},
		`(?i)sendgrid`:                         {Severity: "MEDIUM", Category: "Third Party", Description: "SendGrid integration detected"},
		`(?i)mailchimp`:                        {Severity: "INFO", Category: "Third Party", Description: "Mailchimp integration detected"},
		`(?i)slack[_-]?webhook`:                {Severity: "MEDIUM", Category: "Third Party", Description: "Slack webhook found"},
		`xox[baprs]-[0-9a-zA-Z]{10,}`:          {Severity: "HIGH", Category: "Third Party", Description: "Slack token found"},

		// Database Credentials
		`(?i)mongodb(\+srv)?://`:               {Severity: "HIGH", Category: "Database", Description: "MongoDB connection string found"},
		`(?i)mysql://`:                         {Severity: "HIGH", Category: "Database", Description: "MySQL connection string found"},
		`(?i)postgres(ql)?://`:                 {Severity: "HIGH", Category: "Database", Description: "PostgreSQL connection string found"},
		`(?i)redis://`:                         {Severity: "HIGH", Category: "Database", Description: "Redis connection string found"},
		`(?i)jdbc:`:                            {Severity: "MEDIUM", Category: "Database", Description: "JDBC connection string found"},

		// Insecure Communication (MASTG-TEST-0006)
		`http://[^localhost][^\s"'<>]+`:        {Severity: "MEDIUM", Category: "Insecure Communication", Description: "HTTP URL found (non-HTTPS)"},
		`(?i)ssl[_-]?verify\s*[=:]\s*(false|0|no)`: {Severity: "HIGH", Category: "Insecure Communication", Description: "SSL verification disabled"},
		`(?i)trust[_-]?all[_-]?cert`:           {Severity: "HIGH", Category: "Insecure Communication", Description: "Trust all certificates pattern found"},
		`(?i)allow[_-]?all[_-]?hostname`:       {Severity: "HIGH", Category: "Insecure Communication", Description: "Allow all hostnames pattern found"},
		`(?i)insecure[_-]?ssl`:                 {Severity: "HIGH", Category: "Insecure Communication", Description: "Insecure SSL configuration found"},
		`setHostnameVerifier`:                  {Severity: "MEDIUM", Category: "Insecure Communication", Description: "Custom hostname verifier detected"},
		`TrustManager`:                         {Severity: "MEDIUM", Category: "Insecure Communication", Description: "Custom TrustManager detected"},
		`X509TrustManager`:                     {Severity: "MEDIUM", Category: "Insecure Communication", Description: "Custom X509TrustManager detected"},

		// Cryptography Issues (MASTG-TEST-0013, MASTG-TEST-0014)
		`(?i)DES[^C]`:                          {Severity: "HIGH", Category: "Weak Cryptography", Description: "DES encryption (weak) detected"},
		`(?i)3DES`:                             {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "3DES encryption (deprecated) detected"},
		`(?i)RC4`:                              {Severity: "HIGH", Category: "Weak Cryptography", Description: "RC4 encryption (weak) detected"},
		`(?i)MD5`:                              {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "MD5 hash (weak) detected"},
		`(?i)SHA-?1[^0-9]`:                     {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "SHA-1 hash (deprecated) detected"},
		`(?i)ECB`:                              {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "ECB mode (insecure) detected"},
		`(?i)PKCS1Padding`:                     {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "PKCS1Padding (vulnerable) detected"},
		`(?i)hardcoded[_-]?iv`:                 {Severity: "HIGH", Category: "Weak Cryptography", Description: "Hardcoded IV reference found"},
		`(?i)static[_-]?iv`:                    {Severity: "HIGH", Category: "Weak Cryptography", Description: "Static IV reference found"},
		`SecureRandom`:                         {Severity: "INFO", Category: "Cryptography", Description: "SecureRandom usage detected"},

		// Root/Jailbreak Detection
		`(?i)/system/app/Superuser`:            {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser) found"},
		`(?i)/system/xbin/su`:                  {Severity: "INFO", Category: "Root Detection", Description: "Root detection (su binary) found"},
		`(?i)/sbin/su`:                         {Severity: "INFO", Category: "Root Detection", Description: "Root detection (sbin su) found"},
		`(?i)com\.noshufou\.android\.su`:       {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser app) found"},
		`(?i)com\.thirdparty\.superuser`:       {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser) found"},
		`(?i)eu\.chainfire\.supersu`:           {Severity: "INFO", Category: "Root Detection", Description: "Root detection (SuperSU) found"},
		`(?i)com\.koushikdutta\.superuser`:     {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser) found"},
		`(?i)com\.topjohnwu\.magisk`:           {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Magisk) found"},
		`(?i)isDeviceRooted`:                   {Severity: "INFO", Category: "Root Detection", Description: "Root detection method found"},
		`(?i)checkRoot`:                        {Severity: "INFO", Category: "Root Detection", Description: "Root detection method found"},
		`(?i)RootBeer`:                         {Severity: "INFO", Category: "Root Detection", Description: "RootBeer library detected"},

		// Anti-Tampering & Anti-Debugging
		`(?i)frida`:                            {Severity: "INFO", Category: "Anti-Tampering", Description: "Frida detection code found"},
		`(?i)xposed`:                           {Severity: "INFO", Category: "Anti-Tampering", Description: "Xposed detection code found"},
		`(?i)substrate`:                        {Severity: "INFO", Category: "Anti-Tampering", Description: "Cydia Substrate detection found"},
		`(?i)magisk`:                           {Severity: "INFO", Category: "Anti-Tampering", Description: "Magisk detection code found"},
		`(?i)isDebuggerConnected`:              {Severity: "INFO", Category: "Anti-Debugging", Description: "Debugger detection found"},
		`(?i)android\.os\.Debug`:               {Severity: "INFO", Category: "Anti-Debugging", Description: "Debug class usage found"},
		`(?i)ptrace`:                           {Severity: "INFO", Category: "Anti-Debugging", Description: "Ptrace anti-debugging found"},
		`(?i)TracerPid`:                        {Severity: "INFO", Category: "Anti-Debugging", Description: "TracerPid check found"},

		// Emulator Detection
		`(?i)generic.*Build`:                   {Severity: "INFO", Category: "Emulator Detection", Description: "Emulator detection pattern found"},
		`(?i)goldfish`:                         {Severity: "INFO", Category: "Emulator Detection", Description: "Emulator detection (goldfish) found"},
		`(?i)isEmulator`:                       {Severity: "INFO", Category: "Emulator Detection", Description: "Emulator detection method found"},

		// Logging & Debug Output
		`(?i)Log\.(d|v|i|w|e)\s*\(`:            {Severity: "LOW", Category: "Information Disclosure", Description: "Android logging detected"},
		`(?i)System\.out\.print`:               {Severity: "LOW", Category: "Information Disclosure", Description: "System.out logging detected"},
		`(?i)printStackTrace`:                  {Severity: "LOW", Category: "Information Disclosure", Description: "Stack trace printing detected"},
		`(?i)BuildConfig\.DEBUG`:               {Severity: "INFO", Category: "Debug Code", Description: "BuildConfig.DEBUG check found"},

		// WebView Issues (MASTG-TEST-0031)
		`setJavaScriptEnabled`:                 {Severity: "MEDIUM", Category: "WebView", Description: "JavaScript enabled in WebView"},
		`addJavascriptInterface`:               {Severity: "HIGH", Category: "WebView", Description: "JavaScript interface in WebView (XSS risk)"},
		`setAllowFileAccess`:                   {Severity: "MEDIUM", Category: "WebView", Description: "File access in WebView"},
		`setAllowUniversalAccessFromFileURLs`:  {Severity: "HIGH", Category: "WebView", Description: "Universal file access in WebView"},
		`setAllowFileAccessFromFileURLs`:       {Severity: "MEDIUM", Category: "WebView", Description: "File URL access in WebView"},

		// Data Storage (MASTG-TEST-0001, MASTG-TEST-0002)
		`MODE_WORLD_READABLE`:                  {Severity: "HIGH", Category: "Data Storage", Description: "World-readable file mode found"},
		`MODE_WORLD_WRITEABLE`:                 {Severity: "HIGH", Category: "Data Storage", Description: "World-writable file mode found"},
		`getExternalStorageDirectory`:          {Severity: "MEDIUM", Category: "Data Storage", Description: "External storage usage detected"},
		`getExternalFilesDir`:                  {Severity: "INFO", Category: "Data Storage", Description: "External files directory usage"},
		`SharedPreferences`:                    {Severity: "INFO", Category: "Data Storage", Description: "SharedPreferences usage detected"},

		// Intent/IPC Vulnerabilities
		`(?i)intent\.setData`:                  {Severity: "INFO", Category: "IPC", Description: "Intent data setting detected"},
		`(?i)intent\.putExtra`:                 {Severity: "INFO", Category: "IPC", Description: "Intent extras detected"},
		`(?i)PendingIntent`:                    {Severity: "INFO", Category: "IPC", Description: "PendingIntent usage detected"},
		`(?i)exported\s*=\s*"?true`:            {Severity: "MEDIUM", Category: "IPC", Description: "Exported component detected"},

		// SQL Injection (MASTG-TEST-0026)
		`(?i)rawQuery`:                         {Severity: "MEDIUM", Category: "SQL Injection", Description: "Raw SQL query detected"},
		`(?i)execSQL`:                          {Severity: "MEDIUM", Category: "SQL Injection", Description: "execSQL usage detected"},
		`(?i)SELECT.*FROM.*WHERE`:              {Severity: "LOW", Category: "SQL", Description: "SQL query pattern found"},

		// Backup Indicators
		`(?i)backup[_-]?agent`:                 {Severity: "MEDIUM", Category: "Backup", Description: "Backup agent reference found"},
		`(?i)BackupManager`:                    {Severity: "INFO", Category: "Backup", Description: "BackupManager usage detected"},

		// Sensitive URLs/Endpoints
		`(?i)/api/`:                            {Severity: "INFO", Category: "API Endpoint", Description: "API endpoint path found"},
		`(?i)/v[0-9]+/`:                        {Severity: "INFO", Category: "API Endpoint", Description: "Versioned API path found"},
		`(?i)\.amazonaws\.com`:                 {Severity: "MEDIUM", Category: "Cloud", Description: "AWS endpoint found"},
		`(?i)\.blob\.core\.windows\.net`:       {Severity: "MEDIUM", Category: "Cloud", Description: "Azure blob storage found"},
		`(?i)\.s3\.`:                           {Severity: "MEDIUM", Category: "Cloud", Description: "AWS S3 bucket found"},

		// Biometric/Authentication
		`(?i)BiometricPrompt`:                  {Severity: "INFO", Category: "Authentication", Description: "Biometric authentication detected"},
		`(?i)FingerprintManager`:               {Severity: "INFO", Category: "Authentication", Description: "Fingerprint authentication detected"},
		`(?i)KeyguardManager`:                  {Severity: "INFO", Category: "Authentication", Description: "Keyguard manager usage detected"},
	}

	for pattern, issue := range sensitivePatterns {
		matches, err := SearchStrings(apkPath, pattern)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			issue.Details = fmt.Sprintf("Found in %d file(s)", len(matches))
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// GetDexFiles returns all DEX files in the APK.
func GetDexFiles(apkPath string) ([]string, error) {
	files, err := ListFiles(apkPath)
	if err != nil {
		return nil, err
	}

	var dexFiles []string
	for _, f := range files {
		if strings.HasSuffix(f, ".dex") {
			dexFiles = append(dexFiles, f)
		}
	}
	return dexFiles, nil
}

// GetNativeLibraries returns all native library files in the APK.
func GetNativeLibraries(apkPath string) (map[string][]string, error) {
	files, err := ListFiles(apkPath)
	if err != nil {
		return nil, err
	}

	libs := make(map[string][]string)
	for _, f := range files {
		if strings.HasPrefix(f, "lib/") && strings.HasSuffix(f, ".so") {
			parts := strings.Split(f, "/")
			if len(parts) >= 3 {
				arch := parts[1]
				libs[arch] = append(libs[arch], parts[2])
			}
		}
	}
	return libs, nil
}

// Binary XML parsing helpers

// parseAXML parses Android binary XML format.
func parseAXML(data []byte) (string, error) {
	if len(data) < 8 {
		return "", fmt.Errorf("invalid AXML: too short")
	}

	// Check magic number
	magic := binary.LittleEndian.Uint32(data[0:4])
	if magic != 0x00080003 {
		return "", fmt.Errorf("invalid AXML magic number")
	}

	// For now, return a placeholder - full AXML parsing is complex
	return "[Binary XML - use aapt2 or apktool for full parsing]", nil
}

// ReadManifestRaw reads the raw AndroidManifest.xml bytes.
func ReadManifestRaw(apkPath string) ([]byte, error) {
	zr, err := OpenAPK(apkPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == "AndroidManifest.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("AndroidManifest.xml not found")
}

// ParseManifestBasic extracts basic info from manifest using aapt2 if available.
func ParseManifestBasic(apkPath string) (*APKInfo, error) {
	info, err := GetAPKInfo(apkPath)
	if err != nil {
		return nil, err
	}

	// Try aapt2 first (more reliable), then aapt, then fall back to string extraction
	if perms, pkg := extractPermissionsWithAapt(apkPath); len(perms) > 0 {
		info.Permissions = perms
		if pkg != "" {
			info.PackageName = pkg
		}
		return info, nil
	}

	// Fall back to string extraction from binary manifest
	manifestData, err := ReadManifestRaw(apkPath)
	if err != nil {
		return info, nil
	}

	// Extract readable strings from binary manifest
	extractedStrings := extractStringsFromBytes(manifestData, 4)
	seen := make(map[string]bool)
	for _, s := range extractedStrings {
		// Look for package name pattern
		if strings.Contains(s, ".") && !strings.Contains(s, " ") && len(s) > 5 {
			if isPermissionString(s) && !seen[s] {
				info.Permissions = append(info.Permissions, s)
				seen[s] = true
			} else if info.PackageName == "" && isValidPackageName(s) {
				info.PackageName = s
			}
		}
	}

	return info, nil
}

// extractPermissionsWithAapt uses aapt2 or aapt to extract permissions from APK.
func extractPermissionsWithAapt(apkPath string) ([]string, string) {
	var permissions []string
	var packageName string

	// Try aapt2 first
	tools := []string{"aapt2", "aapt"}
	for _, tool := range tools {
		cmd := exec.Command(tool, "dump", "permissions", apkPath)
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			permissions, packageName = parseAaptPermissionOutput(string(output))
			if len(permissions) > 0 {
				return permissions, packageName
			}
		}

		// Try badging for package name and permissions
		cmd = exec.Command(tool, "dump", "badging", apkPath)
		output, err = cmd.Output()
		if err == nil && len(output) > 0 {
			perms, pkg := parseAaptBadgingOutput(string(output))
			if len(perms) > 0 {
				return perms, pkg
			}
		}
	}

	return nil, ""
}

// parseAaptPermissionOutput parses the output of aapt dump permissions.
func parseAaptPermissionOutput(output string) ([]string, string) {
	var permissions []string
	var packageName string
	seen := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				packageName = parts[1]
			}
		} else if strings.HasPrefix(line, "uses-permission:") {
			// Format: uses-permission: name='android.permission.INTERNET'
			if idx := strings.Index(line, "name='"); idx != -1 {
				start := idx + 6
				end := strings.Index(line[start:], "'")
				if end != -1 {
					perm := line[start : start+end]
					if !seen[perm] {
						permissions = append(permissions, perm)
						seen[perm] = true
					}
				}
			} else {
				// Alternative format: uses-permission: android.permission.INTERNET
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					perm := parts[1]
					if !seen[perm] {
						permissions = append(permissions, perm)
						seen[perm] = true
					}
				}
			}
		}
	}

	return permissions, packageName
}

// parseAaptBadgingOutput parses the output of aapt dump badging.
func parseAaptBadgingOutput(output string) ([]string, string) {
	var permissions []string
	var packageName string
	seen := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package:") {
			// Format: package: name='com.example.app' versionCode='1' ...
			if idx := strings.Index(line, "name='"); idx != -1 {
				start := idx + 6
				end := strings.Index(line[start:], "'")
				if end != -1 {
					packageName = line[start : start+end]
				}
			}
		} else if strings.HasPrefix(line, "uses-permission:") {
			// Format: uses-permission: name='android.permission.INTERNET'
			if idx := strings.Index(line, "name='"); idx != -1 {
				start := idx + 6
				end := strings.Index(line[start:], "'")
				if end != -1 {
					perm := line[start : start+end]
					if !seen[perm] {
						permissions = append(permissions, perm)
						seen[perm] = true
					}
				}
			}
		}
	}

	return permissions, packageName
}

// isPermissionString checks if a string looks like an Android permission.
func isPermissionString(s string) bool {
	// Check common permission prefixes
	permissionPrefixes := []string{
		"android.permission.",
		"com.google.android.",
		"com.android.",
	}

	for _, prefix := range permissionPrefixes {
		if strings.HasPrefix(s, prefix) {
			// Verify it looks like a permission (contains "permission" or is uppercase after prefix)
			remainder := strings.TrimPrefix(s, prefix)
			if strings.Contains(s, "permission") || isUpperSnakeCase(remainder) {
				return true
			}
		}
	}

	// Check for custom permissions that contain ".permission."
	if strings.Contains(s, ".permission.") {
		return true
	}

	return false
}

// isUpperSnakeCase checks if string is UPPER_SNAKE_CASE (common for permissions).
func isUpperSnakeCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

// isValidPackageName checks if a string looks like a valid package name.
func isValidPackageName(s string) bool {
	if len(s) < 3 || len(s) > 200 {
		return false
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
		// Check if part starts with letter and contains only valid chars
		for i, c := range part {
			if i == 0 {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_') {
					return false
				}
			} else {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
					return false
				}
			}
		}
	}
	return true
}

// HasFile checks if a file exists in the APK.
func HasFile(apkPath, fileName string) (bool, error) {
	zr, err := OpenAPK(apkPath)
	if err != nil {
		return false, err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == fileName {
			return true, nil
		}
	}
	return false, nil
}

// GetFileSize returns the size of a file in the APK.
func GetFileSize(apkPath, fileName string) (uint64, error) {
	zr, err := OpenAPK(apkPath)
	if err != nil {
		return 0, err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == fileName {
			return f.UncompressedSize64, nil
		}
	}
	return 0, fmt.Errorf("file not found: %s", fileName)
}

// XML structures for parsing decompiled manifest

// Manifest represents the AndroidManifest.xml structure.
type Manifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	Package     string      `xml:"package,attr"`
	VersionCode string      `xml:"versionCode,attr"`
	VersionName string      `xml:"versionName,attr"`
	UsesSDK     UsesSDK     `xml:"uses-sdk"`
	Permissions []UsesPermission `xml:"uses-permission"`
	Application Application `xml:"application"`
}

// UsesSDK represents SDK version requirements.
type UsesSDK struct {
	MinSDK    string `xml:"minSdkVersion,attr"`
	TargetSDK string `xml:"targetSdkVersion,attr"`
}

// UsesPermission represents a permission declaration.
type UsesPermission struct {
	Name string `xml:"name,attr"`
}

// Application represents the application element.
type Application struct {
	Debuggable  string     `xml:"debuggable,attr"`
	AllowBackup string     `xml:"allowBackup,attr"`
	Label       string     `xml:"label,attr"`
	Activities  []Activity `xml:"activity"`
	Services    []Service  `xml:"service"`
	Receivers   []Receiver `xml:"receiver"`
	Providers   []Provider `xml:"provider"`
}

// Activity represents an activity declaration.
type Activity struct {
	Name     string `xml:"name,attr"`
	Exported string `xml:"exported,attr"`
}

// Service represents a service declaration.
type Service struct {
	Name     string `xml:"name,attr"`
	Exported string `xml:"exported,attr"`
}

// Receiver represents a broadcast receiver declaration.
type Receiver struct {
	Name     string `xml:"name,attr"`
	Exported string `xml:"exported,attr"`
}

// Provider represents a content provider declaration.
type Provider struct {
	Name        string `xml:"name,attr"`
	Exported    string `xml:"exported,attr"`
	Authorities string `xml:"authorities,attr"`
}

// ParseXMLManifest parses a decompiled XML manifest.
func ParseXMLManifest(xmlData []byte) (*Manifest, error) {
	var manifest Manifest
	err := xml.Unmarshal(xmlData, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest XML: %v", err)
	}
	return &manifest, nil
}

// FormatSize formats a file size in human-readable format.
func FormatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// StringContains is a helper to avoid import issues with strings package.
func stringContains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

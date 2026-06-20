package apk

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
	"android.permission.READ_PHONE_STATE":       {"READ_PHONE_STATE", "Read phone state and identity (IMEI)", "Privacy"},
	"android.permission.READ_PHONE_NUMBERS":     {"READ_PHONE_NUMBERS", "Read phone numbers", "Privacy"},
	"android.permission.CALL_PHONE":             {"CALL_PHONE", "Make phone calls without user action", "Financial"},
	"android.permission.ANSWER_PHONE_CALLS":     {"ANSWER_PHONE_CALLS", "Answer incoming calls programmatically", "Privacy"},
	"android.permission.READ_CALL_LOG":          {"READ_CALL_LOG", "Read call history", "Privacy"},
	"android.permission.WRITE_CALL_LOG":         {"WRITE_CALL_LOG", "Modify call history", "Privacy"},
	"android.permission.ADD_VOICEMAIL":          {"ADD_VOICEMAIL", "Add voicemail messages", "Privacy"},
	"android.permission.USE_SIP":                {"USE_SIP", "Use SIP service for calls", "Financial"},
	"android.permission.PROCESS_OUTGOING_CALLS": {"PROCESS_OUTGOING_CALLS", "Intercept outgoing calls", "Privacy"},
	"android.permission.ACCEPT_HANDOVER":        {"ACCEPT_HANDOVER", "Accept call handover from another app", "Privacy"},

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

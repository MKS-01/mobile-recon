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
var DangerousPermissions = map[string]DangerousPermission{
	"android.permission.READ_CALENDAR":          {"READ_CALENDAR", "Read calendar events", "Privacy"},
	"android.permission.WRITE_CALENDAR":         {"WRITE_CALENDAR", "Modify calendar events", "Privacy"},
	"android.permission.CAMERA":                 {"CAMERA", "Access device camera", "Privacy"},
	"android.permission.READ_CONTACTS":          {"READ_CONTACTS", "Read user contacts", "Privacy"},
	"android.permission.WRITE_CONTACTS":         {"WRITE_CONTACTS", "Modify user contacts", "Privacy"},
	"android.permission.GET_ACCOUNTS":           {"GET_ACCOUNTS", "Access device accounts", "Privacy"},
	"android.permission.ACCESS_FINE_LOCATION":   {"ACCESS_FINE_LOCATION", "Access precise location", "Privacy/Tracking"},
	"android.permission.ACCESS_COARSE_LOCATION": {"ACCESS_COARSE_LOCATION", "Access approximate location", "Privacy/Tracking"},
	"android.permission.RECORD_AUDIO":           {"RECORD_AUDIO", "Record audio from microphone", "Privacy"},
	"android.permission.READ_PHONE_STATE":       {"READ_PHONE_STATE", "Read phone state and identity", "Privacy"},
	"android.permission.READ_PHONE_NUMBERS":     {"READ_PHONE_NUMBERS", "Read phone numbers", "Privacy"},
	"android.permission.CALL_PHONE":             {"CALL_PHONE", "Make phone calls", "Financial"},
	"android.permission.ANSWER_PHONE_CALLS":     {"ANSWER_PHONE_CALLS", "Answer incoming calls", "Privacy"},
	"android.permission.READ_CALL_LOG":          {"READ_CALL_LOG", "Read call history", "Privacy"},
	"android.permission.WRITE_CALL_LOG":         {"WRITE_CALL_LOG", "Modify call history", "Privacy"},
	"android.permission.ADD_VOICEMAIL":          {"ADD_VOICEMAIL", "Add voicemail messages", "Privacy"},
	"android.permission.USE_SIP":                {"USE_SIP", "Use SIP service", "Financial"},
	"android.permission.BODY_SENSORS":           {"BODY_SENSORS", "Access body sensors", "Health/Privacy"},
	"android.permission.SEND_SMS":               {"SEND_SMS", "Send SMS messages", "Financial"},
	"android.permission.RECEIVE_SMS":            {"RECEIVE_SMS", "Receive SMS messages", "Privacy"},
	"android.permission.READ_SMS":               {"READ_SMS", "Read SMS messages", "Privacy"},
	"android.permission.RECEIVE_WAP_PUSH":       {"RECEIVE_WAP_PUSH", "Receive WAP push messages", "Privacy"},
	"android.permission.RECEIVE_MMS":            {"RECEIVE_MMS", "Receive MMS messages", "Privacy"},
	"android.permission.READ_EXTERNAL_STORAGE":  {"READ_EXTERNAL_STORAGE", "Read external storage", "Privacy"},
	"android.permission.WRITE_EXTERNAL_STORAGE": {"WRITE_EXTERNAL_STORAGE", "Write to external storage", "Data"},
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

	// Search for sensitive patterns
	sensitivePatterns := map[string]SecurityIssue{
		`(?i)api[_-]?key`:           {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential API key found"},
		`(?i)secret`:                {Severity: "MEDIUM", Category: "Hardcoded Secrets", Description: "Potential secret found"},
		`(?i)password`:              {Severity: "MEDIUM", Category: "Hardcoded Secrets", Description: "Potential password reference found"},
		`(?i)private[_-]?key`:       {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential private key found"},
		`(?i)aws[_-]?access`:        {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential AWS credentials found"},
		`(?i)firebase`:              {Severity: "INFO", Category: "Third Party", Description: "Firebase integration detected"},
		`http://`:                   {Severity: "MEDIUM", Category: "Insecure Communication", Description: "HTTP URL found (non-HTTPS)"},
		`(?i)root`:                  {Severity: "INFO", Category: "Root Detection", Description: "Root-related string found"},
		`(?i)su`:                    {Severity: "INFO", Category: "Root Detection", Description: "Superuser reference found"},
		`(?i)frida`:                 {Severity: "INFO", Category: "Anti-Tampering", Description: "Frida detection code found"},
		`(?i)xposed`:                {Severity: "INFO", Category: "Anti-Tampering", Description: "Xposed detection code found"},
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

	// Try to extract strings that look like package names and permissions
	manifestData, err := ReadManifestRaw(apkPath)
	if err != nil {
		return info, nil
	}

	// Extract readable strings from binary manifest
	extractedStrings := extractStringsFromBytes(manifestData, 4)
	for _, s := range extractedStrings {
		// Look for package name pattern
		if strings.Contains(s, ".") && !strings.Contains(s, " ") && len(s) > 5 {
			if strings.HasPrefix(s, "android.permission.") {
				info.Permissions = append(info.Permissions, s)
			} else if info.PackageName == "" && isValidPackageName(s) {
				info.PackageName = s
			}
		}
	}

	return info, nil
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

package apk

import (
	"encoding/xml"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// ManifestAttribute represents an attribute in the Android manifest.
type ManifestAttribute struct {
	Name  string
	Value string
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

// Manifest represents the AndroidManifest.xml structure.
type Manifest struct {
	XMLName     xml.Name         `xml:"manifest"`
	Package     string           `xml:"package,attr"`
	VersionCode string           `xml:"versionCode,attr"`
	VersionName string           `xml:"versionName,attr"`
	UsesSDK     UsesSDK          `xml:"uses-sdk"`
	Permissions []UsesPermission `xml:"uses-permission"`
	Application Application      `xml:"application"`
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

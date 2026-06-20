// Package apk provides core functionality for analyzing Android APK files.
// It includes APK parsing, manifest extraction, and security analysis.
package apk

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"
)

// APKInfo contains metadata extracted from an APK file.
type APKInfo struct {
	FilePath      string
	FileSize      int64
	PackageName   string
	VersionName   string
	VersionCode   int
	MinSDK        int
	TargetSDK     int
	Permissions   []string
	Activities    []string
	Services      []string
	Receivers     []string
	Providers     []string
	HasNativeLib  bool
	Architectures []string
	IsDebuggable  bool
	AllowBackup   bool
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

package apk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

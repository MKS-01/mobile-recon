package apk

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

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

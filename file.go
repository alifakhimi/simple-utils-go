package simutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateFileOption holds optional configuration for creating files and directories
type CreateFileOption struct {
	// Directory permission (e.g., 755 for rwxr-xr-x)
	DirPerm os.FileMode
	// File permission (e.g., 644 for rw-r--r--)
	FilePerm os.FileMode
	// File flag (e.g., os.O_WRONLY|os.O_CREATE|os.O_TRUNC for write-only, create if not exist, and truncate)
	FileFlag int
}

// FileNameFriendlyNowTime generates a file-system-friendly string representing the current time
func FileNameFriendlyNowTime() string {
	return FileNameFriendlyTime(time.Now().Local())
}

// FileNameFriendlyTime converts a time to a string without characters invalid for file names
func FileNameFriendlyTime(t time.Time) string {
	// Format the time as an ISO 8601 string (e.g., "2025-01-07T15:04:05Z07:00")
	ts := t.Format(time.RFC3339)
	// Replace invalid characters (':' and '-') with nothing to make it file-system-friendly
	return strings.Replace(strings.Replace(ts, ":", "", -1), "-", "", -1)
}

// CreateFile ensures the directory exists and opens a file ready for writing
func CreateFile(path string, opt ...CreateFileOption) (*os.File, error) {
	var (
		dirPerm  = os.FileMode(0755)                      // Default directory permissions
		filePerm = os.FileMode(0644)                      // Default file permissions
		fileFlag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC // Default file flags
	)

	// Override defaults with provided options, if any
	for _, o := range opt {
		if o.DirPerm > 0 {
			dirPerm = o.DirPerm
		}
		if o.FilePerm > 0 {
			filePerm = o.FilePerm
		}
		if o.FileFlag > 0 {
			fileFlag = o.FileFlag
		}
	}

	// Extract the directory part of the provided file path
	dir := filepath.Dir(path)

	// Ensure the directory exists (creates intermediate directories if needed)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	// Open or create the file with the specified flags and permissions
	file, err := os.OpenFile(path, fileFlag, filePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create or open file: %w", err)
	}

	return file, nil
}

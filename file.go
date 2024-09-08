package simutils

import (
	"fmt"
	"os"
	"path/filepath"
)

type CreateFileOption struct {
	// e.g. 755
	DirPerm os.FileMode
	// e.g. 644
	FilePerm os.FileMode
	// e.g. os.O_WRONLY|os.O_CREATE|os.O_TRUNC
	FileFlag int
}

// CreateFile ensures that the directory exists and returns an open file ready for writing
func CreateFile(path string, opt ...CreateFileOption) (*os.File, error) {
	var (
		dirPerm  = os.FileMode(0755)
		filePerm = os.FileMode(0644)
		fileFlag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	)

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

	// Get the directory path from the file path
	dir := filepath.Dir(path)

	// Create the directory along the path if it doesn't exist
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	// Open the file, create if it doesn't exist and truncate it if it does
	file, err := os.OpenFile(path, fileFlag, filePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create or open file: %w", err)
	}

	return file, nil
}

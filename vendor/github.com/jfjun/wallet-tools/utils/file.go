package utils

import "os"

// OpenFile opens or creates a file
// If the file already exists, open it . If it does not,
// It will create the file with mode 0644.
func OpenFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return f, err
}

// FileExist determines whether a file exists
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

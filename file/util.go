package file

import (
	"os"
	"path/filepath"
)

// IsExists checks if a file or directory exists.
func IsExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir checks if the given path is a directory.
func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetFileList returns a slice of file paths for all files in the given directory and its subdirectories.
func GetFileList(dir string) ([]string, error) {
	var res []string

	if !IsDir(dir) {
		exists, err := IsExists(dir)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, nil
		}
		res = append(res, dir)
		return res, nil
	}

	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fi := range fileList {
		fullPath := filepath.Join(dir, fi.Name())
		if fi.IsDir() {
			subFiles, err := GetFileList(fullPath)
			if err != nil {
				return nil, err
			}
			res = append(res, subFiles...)
		}
		res = append(res, fullPath)
	}

	return res, nil
}

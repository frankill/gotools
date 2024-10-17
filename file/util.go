package file

import (
	"os"
	"path/filepath"
)

// 判断文件或文件夹是否存在
// 如果是，返回 true，否则返回 false
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

// 判断所有文件夹是否存在
func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// 获取文件夹下所有文件
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

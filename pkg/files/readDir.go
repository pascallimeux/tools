package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAbsDirName(searchDir string) (string, error) {
	if stat, err := os.Stat(searchDir); err == nil && stat.IsDir() {
		absSearchDir, err := filepath.Abs(searchDir)
		if err != nil {
			return "", fmt.Errorf("%s", err.Error())
		}
		return absSearchDir, nil
	}
	return "", fmt.Errorf("%s is not a directory", searchDir)
}

func ReadDir(searchDir string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		info, err := os.Stat(path)
		if err == nil && !info.IsDir(){
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList, err
}

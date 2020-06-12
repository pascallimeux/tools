package files

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// CopyFile copy a file with is absolute path in dstFolder
func CopyFile(src, dstFolder string) (int64, error) {

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	LastIndex := strings.LastIndex(src, "/")
	fileFolder := src[:LastIndex]
	if fileFolder != "" {
		os.MkdirAll(dstFolder+"/"+fileFolder, os.ModePerm)
	}
	dst := dstFolder + "/" + src

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

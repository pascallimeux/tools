package main

import (
	"fmt"
	"os"
	"tools/pkg/files"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [path]\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	sdir := os.Args[1]
	//sdir := "/Users/pascallimeux/Library/Application Support/Google/Chrome"
	searchDir, err := files.GetAbsDirName(sdir)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}

	fileList, err := files.ReadDir(searchDir)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}

	sqlLiteFileList := files.GetFileList(fileList, files.SqliteSign)

	for _, filename := range sqlLiteFileList {
		size, err := files.CopyFile(filename, "/tmp/test/")
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(-1)
		}
		fmt.Printf("copy file %s (%d)\n", filename, size)
	}
}
// ./getSqlLiteDB "/Users/pascallimeux/Library/Application Support/Google/Chrome"

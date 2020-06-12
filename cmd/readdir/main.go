// read recursiv dir
// go build -o readdir main.go

package main

import (
	"fmt"
	"os"
	"tools/pkg/files"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %v path...\n", os.Args[0])
		os.Exit(0)
	}

	searchDir, err := files.GetAbsDirName(os.Args[1])
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}

	fileList, err := files.ReadDir(searchDir)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}

	for _, file := range fileList {
		fmt.Println(file)
	}
}

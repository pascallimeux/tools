// forensic file use signature file

package main

import (
	"fmt"
	"tools/pkg/files"
)

func main() {
	filename := "pngfile"
	format, err := files.GetFileFormat(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s has format:%s\n", filename, format)

	filename = "jpgfile"
	format, err = files.GetFileFormat(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s has format:%s\n", filename, format)
}

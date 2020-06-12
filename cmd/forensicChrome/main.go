package main

import (
	"fmt"
	"tools/internal/forensicChrome"
)

func main() {
	err := forensicChrome.Process("/tmp/forensic")
	if err != nil {
		fmt.Println(err.Error())
	}
}

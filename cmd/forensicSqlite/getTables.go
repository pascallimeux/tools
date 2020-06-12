package main

import (
	"tools/pkg/sqllite3"
)

func main() {
	//if len(os.Args) == 1 {
	//	fmt.Printf("usage: %v path...\n", os.Args[0])
	//	os.Exit(0)
	//}
	//sqllite3.GetDBTables(os.Args[1])
	sqllite3.GetDBTables("/Users/pascallimeux/Library/Application Support/Google/Chrome/Users/pascallimeux/Library/Application Support/Google/Chrome/Default/History")
	sqllite3.GetColumns("/Users/pascallimeux/Library/Application Support/Google/Chrome/Users/pascallimeux/Library/Application Support/Google/Chrome/Default/History", "urls")
	sqllite3.GetTable("/Users/pascallimeux/Library/Application Support/Google/Chrome/Users/pascallimeux/Library/Application Support/Google/Chrome/Default/History", "urls")

}

// ./getTables "/Users/pascallimeux/Library/Application Support/Google/Chrome/Users/pascallimeux/Library/Application Support/Google/Chrome/Default/History"

package sqllite3

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	name string
)

func GetDBTables(dbname string) {
	db, err := sql.Open("sqlite3", dbname)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type ='table'")
	checkErr(err)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		checkErr(err)
		log.Println(name)
	}
	err = rows.Err()
	checkErr(err)
}

func GetColumns(dbname, tableName string) {
	db, err := sql.Open("sqlite3", dbname)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE tbl_name = '%s' AND type = 'table'", tableName))
	checkErr(err)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		checkErr(err)
		log.Println(name)
	}
	err = rows.Err()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

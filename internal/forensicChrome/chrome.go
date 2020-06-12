// +build darwin

package forensicChrome

import (
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"time"
	"tools/pkg/files"

	_ "github.com/mattn/go-sqlite3"
)

type CsvStruct interface {
	GetHeaders() []string
	ToSlice() []string
}

type Login struct {
	Origin_url     string
	Username_value string
	Password_value []byte
}

func (l Login) GetHeaders() []string {
	return []string{"url", "username", "password"}
}

func (l Login) ToSlice() []string {
	return []string{l.Origin_url, l.Username_value, hex.EncodeToString(l.Password_value)}
}

type Term struct {
	Value string
}

func (t Term) GetHeaders() []string {
	return []string{"term"}
}
func (t Term) ToSlice() []string {
	return []string{t.Value}
}

type Download struct {
	Url string
}

func (d Download) GetHeaders() []string {
	return []string{"url"}
}

func (d Download) ToSlice() []string {
	return []string{d.Url}
}

type History struct {
	Url             string
	Title           string
	Visit_count     int
	Last_visit_time time.Time
}

func (h History) GetHeaders() []string {
	return []string{"url", "title", "visit_count", "last_visit_time"}
}

func (h History) ToSlice() []string {
	return []string{h.Url, h.Title, strconv.Itoa(h.Visit_count), h.Last_visit_time.Format("2006-01-02 15:04:05")}
}

func GetChromeLogins(db *sql.DB) (logins []CsvStruct) {
	rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var origin_url, username_value string
		var password_value []byte
		rows.Scan(&origin_url, &username_value, &password_value)
		logins = append(logins, Login{origin_url, username_value, password_value})
	}
	return
}

func GetChromeSearchTerms(db *sql.DB) (terms []CsvStruct) {
	rows, err := db.Query("SELECT term FROM keyword_search_terms")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var term string
		rows.Scan(&term)
		terms = append(terms, Term{term})
	}
	return
}

func GetChromeDownload(db *sql.DB) (downloads []CsvStruct) {
	rows, err := db.Query("SELECT url FROM downloads_url_chains")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var url string
		rows.Scan(&url)
		downloads = append(downloads, Download{url})
	}
	return
}

func GetChromeHistory(db *sql.DB) (histories []CsvStruct) {
	rows, err := db.Query("SELECT url, title, visit_count, last_visit_time FROM urls")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var url, title string
		var visit_count int
		var last_visit_time time.Time
		rows.Scan(&url, &title, &visit_count, &last_visit_time)
		histories = append(histories, History{url, title, visit_count, last_visit_time})
	}
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Process(tmpFolder string) error {
	user, _ := user.Current()
	loginDB := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/Login Data", user.HomeDir)
	historyDB := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/History", user.HomeDir)
	tmpLoginDB := fmt.Sprintf("%s/%s", tmpFolder, loginDB)
	tmpHistoryDB := fmt.Sprintf("%s/%s", tmpFolder, historyDB)

	// copy dbFiles
	_, err := files.CopyFile(loginDB, tmpFolder)
	if err != nil {
		return err
	}
	_, err = files.CopyFile(historyDB, tmpFolder)
	if err != nil {
		return err
	}

	// read dbFiles
	db1, err := sql.Open("sqlite3", tmpLoginDB)
	checkErr(err)
	defer db1.Close()
	logins := GetChromeLogins(db1)

	db2, err := sql.Open("sqlite3", tmpHistoryDB)
	checkErr(err)
	defer db2.Close()
	searchTerms := GetChromeSearchTerms(db2)
	downloads := GetChromeDownload(db2)
	histories := GetChromeHistory(db2)

	// create csvFiles
	err = writeCsv(tmpFolder+"/chromeLogin.csv", logins)
	if err != nil {
		return err
	}
	err = writeCsv(tmpFolder+"/chromeSearchTerms.csv", searchTerms)
	if err != nil {
		return err
	}
	err = writeCsv(tmpFolder+"/chromeDownloads.csv", downloads)
	if err != nil {
		return err
	}
	err = writeCsv(tmpFolder+"/chromeHistories.csv", histories)
	return err
}

func writeCsv(name string, data []CsvStruct) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to wwrite csv file")
	}
	file, err := os.Create(name)
	checkErr(err)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	headers := data[0].GetHeaders()
	if err := w.Write(headers); err != nil {
		return err
	}
	for _, myStruct := range data {
		values := myStruct.ToSlice()
		if err := w.Write(values); err != nil {
			return err
		}
	}
	return nil
}

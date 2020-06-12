// +build darwin

package main

import (
	//"code.google.com/p/go.crypto/pbkdf2"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	_ "github.com/mattn/go-sqlite3"
)

var (
	salt       = "saltysalt"
	iv         = "                "
	length     = 16
	password   = ""
	iterations = 1003
)

type Credentials struct {
	Origin_url     string
	Username_value string
	Password_value []byte
}

func (c *Credentials) DecryptedValue() string {
	if len(c.Password_value) > 0 {
		encryptedValue := c.Password_value[3:]
		return decryptValue(encryptedValue)
	}
	return ""
}

func copyFileToDirectory(pathSourceFile string, pathDestFile string) error {
	sourceFile, err := os.Open(pathSourceFile)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(pathDestFile)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFileInfo, err := destFile.Stat()
	if err != nil {
		return err
	}

	if sourceFileInfo.Size() == destFileInfo.Size() {
	} else {
		return err
	}
	return nil
}

func main() {

	password = "wTFBk3Cs+m2c7NkSI2w9ZA==" //getPassword()

	for _, getCredentials := range getCredentials() {
		fmt.Printf("%s/%s: %s\n", getCredentials.Origin_url, getCredentials.Username_value, getCredentials.DecryptedValue())
	}
}

func decryptValue(encryptedValue []byte) string {
	key := pbkdf2.Key([]byte(password), []byte(salt), iterations, length, sha1.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	decrypted := make([]byte, len(encryptedValue))
	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	cbc.CryptBlocks(decrypted, encryptedValue)

	plainText, err := aesStripPadding(decrypted)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return ""
	}
	return string(plainText)
}

// In the padding scheme the last <padding length> bytes
// have a value equal to the padding length, always in (1,16]
func aesStripPadding(data []byte) ([]byte, error) {
	if len(data)%length != 0 {
		return nil, fmt.Errorf("decrypted data block length is not a multiple of %d", length)
	}
	paddingLen := int(data[len(data)-1])
	if paddingLen > 16 {
		return nil, fmt.Errorf("invalid last block padding length: %d", paddingLen)
	}
	return data[:len(data)-paddingLen], nil
}

func getPassword() string {
	parts := strings.Fields("security find-generic-password -wga Chrome")
	cmd := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(cmd, parts...).Output()
	if err != nil {
		log.Fatal("error finding password ", err)
	}

	return strings.Trim(string(out), "\n")
}

func getCredentials() (credentials []Credentials) {
	usr, _ := user.Current()
	loginFile := fmt.Sprintf("/tmp/test/%s/Library/Application Support/Google/Chrome/Default/Login Data", usr.HomeDir)

	db, err := sql.Open("sqlite3", loginFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var origin_url, username_value string
		var password_value []byte
		rows.Scan(&origin_url, &username_value, &password_value)
		credentials = append(credentials, Credentials{origin_url, username_value, password_value})
	}

	return
}

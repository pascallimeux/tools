package crypto

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func ReadCredentials(username, hostname string) (string, string) {
	var passwordPrompt string

	if username == "" {
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("%s username: ", hostname)
		username, _ = reader.ReadString('\n')
		passwordPrompt = fmt.Sprintf("%s password: ", hostname)
	} else {
		passwordPrompt = fmt.Sprintf("%s@%s password: ", username, hostname)
	}

	fmt.Printf("%s", passwordPrompt)
	bytePassword, _ := terminal.ReadPassword(0)
	password := string(bytePassword)
	fmt.Println("")

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func GenerateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func GenerateEncodedHash(password string) (string, error) {
	hash, err := GenerateHash(password)
	if err != nil {
		return "", err
	}
	return EncodeBase64(hash), nil
}

func CompareHashAndPassword(hashFromDatabase []byte, password string) bool {
	if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(password)); err != nil {
		return false
	}
	return true
}

//CreateAESKeyFromPassword create AES key from plaintext password
func CreateAESKeyFromPassword(plaintextPwd string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(plaintextPwd))
	key := hex.EncodeToString(hasher.Sum(nil))
	//fmt.Println(key)
	return []byte(key)
}

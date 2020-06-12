package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

//AESCipher AES cipher to encrypt and decrypt
type AESCipher struct {
	key []byte
}

// NewAESCipher Instanciate an AES cipher and generate key
func NewAESCipher() (*AESCipher, error) {
	cipher := AESCipher{}
	err := cipher.generateAESKey()
	return &cipher, err
}

func NewAESCipherWithUserPwd(userPwd string) *AESCipher {
	key := CreateAESKeyFromPassword(userPwd)
	cipher := NewAESCipherWithKey(key)
	return cipher
}

// NewAESCipherWithKey Instanciate an AES cipher from base64 encoded key pass in parameter
func NewAESCipherWithEncodedKey(b64Key string) (*AESCipher, error) {
	var err error
	cipher := AESCipher{}
	cipher.key, err = DecodeBase64(b64Key)
	return &cipher, err
}

// NewAESCipherWithKey Instanciate an AES cipher from key pass in parameter
func NewAESCipherWithKey(key []byte) *AESCipher {
	cipher := AESCipher{key: key}
	return &cipher
}

//GetEncodedKey Get base64 encoded key used by the AES cypher
func (a *AESCipher) GetEncodedKey() string {
	encodedB64Key := EncodeBase64(a.key)
	return encodedB64Key
}

func (a *AESCipher) GetKey() []byte {
	return a.key
}

func (a *AESCipher) generateAESKey() error {
	a.key = make([]byte, 32)
	_, err := rand.Read(a.key)
	return err
}

func (a *AESCipher) Encrypt(message []byte) ([]byte, error) {
	start := time.Now()
	c, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("Encrypt AES in %s", time.Since(start)))
	encodedBytes := gcm.Seal(nonce, nonce, message, nil)
	return encodedBytes, nil
}

// EncryptToEncoded byte message with the AES cipher to base64 encoded
// []byte --> crypted []byte --> base 64 encoded
func (a *AESCipher) EncryptToEncoded(message []byte) (string, error) {
	encodedBytes, err := a.Encrypt(message)
	if err != nil {
		return "", err
	}
	encodedB64Text := EncodeBase64(encodedBytes)
	return encodedB64Text, nil
}

func (a *AESCipher) Decrypt(cipherMessage []byte) ([]byte, error) {
	start := time.Now()
	c, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherMessage) < nonceSize {
		return nil, errors.New("cipherMessage too short")
	}

	nonce, cipherMessage := cipherMessage[:nonceSize], cipherMessage[nonceSize:]
	log.Info(fmt.Sprintf("Decrypt AES in %s", time.Since(start)))
	return gcm.Open(nil, nonce, cipherMessage, nil)
}

// Decrypt base 64 encoded message with AES cipher in byte
// base64 encoded --> crypted []byte --> []byte
func (a *AESCipher) DecryptEncodedMessage(encodedB64Text string) ([]byte, error) {
	ciphertext, err := DecodeBase64(encodedB64Text)
	if err != nil {
		return nil, err
	}
	return a.Decrypt(ciphertext)
}

func (a *AESCipher) EncryptFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", filePath)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a folder", filePath)
	}
	fileTime := fileInfo.ModTime()

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	cypherContent, err := a.Encrypt(content)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	err = ioutil.WriteFile(filePath, cypherContent, 0644)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	err = os.Chtimes(filePath, fileTime, fileTime)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return nil
}

func (a *AESCipher) DecryptFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", filePath)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a folder", filePath)
	}
	fileTime := fileInfo.ModTime()

	cypherContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	content, err := a.Decrypt(cypherContent)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	err = os.Chtimes(filePath, fileTime, fileTime)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return nil
}

func (a *AESCipher) EncryptFolder(folderPath string) (int, error) {
	fileInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("%s does not exist", folderPath)
	}
	if !fileInfo.IsDir() {
		return 0, fmt.Errorf("%s is not a folder", folderPath)
	}
	errMsg := ""
	nbFiles := 0
	err = filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				err = a.EncryptFile(path)
				if err != nil {
					errMsg += fmt.Sprintf("error file %s\n", path)
				} else {
					nbFiles += 1
				}
			}
			return nil
		})

	if errMsg != "" {
		return nbFiles, fmt.Errorf("%s", errMsg)
	}
	return nbFiles, err
}

func (a *AESCipher) DecryptFolder(folderPath string) (int, error) {
	fileInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("%s does not exist", folderPath)
	}
	if !fileInfo.IsDir() {
		return 0, fmt.Errorf("%s is not a folder", folderPath)
	}
	errMsg := ""
	nbFiles := 0
	err = filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				err = a.DecryptFile(path)
				if err != nil {
					errMsg += fmt.Sprintf("error file %s\n", path)
				} else {
					nbFiles += 1
				}
			}
			return nil
		})
	if errMsg != "" {
		return nbFiles, fmt.Errorf("%s", errMsg)
	}
	return nbFiles, err
}

func (a *AESCipher) DecryptFileFolder(path string) (int, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("%s does not exist", path)
	}
	if fileInfo.IsDir() {
		return a.DecryptFolder(path)
	} else {
		return 1, a.DecryptFile(path)
	}
}

func (a *AESCipher) EncryptFileFolder(path string) (int, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("%s does not exist", path)
	}
	if fileInfo.IsDir() {
		return a.EncryptFolder(path)
	} else {
		return 1, a.EncryptFile(path)
	}
}

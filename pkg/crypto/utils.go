package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

//GenerateNounceStr generate a string nounce with a specific len
func GenerateNounceStr(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	return str[:len]
}

//GenerateUUID genreate a string UUID (16)
func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10])
	return uuid, nil
}

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GenerateUUIDOld() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	return string(out), err
}

// EncryptMessage encrypts a message with a generated AES key, and encrypts this key with a public RSA key. returns the encrypted key and the encrypted message
func EncryptMessage(message, pubRSAKey []byte) ([]byte, []byte, error) {
	AEScipher, err := NewAESCipher()
	if err != nil {
		return nil, nil, err
	}
	cipherContent, err := AEScipher.Encrypt(message)
	if err != nil {
		return nil, nil, err
	}
	key := AEScipher.GetKey()
	RSACypher, err := NewRSACipherWithKey(nil, pubRSAKey)
	if err != nil {
		return nil, nil, err
	}
	cipherKey, err := RSACypher.Encrypt(key)
	if err != nil {
		return nil, nil, err
	}
	return cipherKey, cipherContent, nil
}

func DecryptMessage(cipherMsg, cipherKey, privRSAKey []byte) ([]byte, error) {
	RSACypher, err := NewRSACipherWithKey(privRSAKey, nil)
	if err != nil {
		return nil, err
	}
	AESKey, err := RSACypher.Decrypt(cipherKey)
	if err != nil {
		return nil, err
	}
	AEScipher := NewAESCipherWithKey(AESKey)
	return AEScipher.Decrypt(cipherMsg)
}

func bytes2Hexa(b []byte) string {
	h := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(h, b)
	//return fmt.Sprintf("0x%s", h)
	return fmt.Sprintf("%s", h)
}

func hexa2Bytes(h string) ([]byte, error) {
	if strings.HasPrefix(h, "0x") {
		h = h[2:]
	}
	return hex.DecodeString(h)
}

func GenerateUUID20HEX() (string, error) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x", b)
	return uuid, nil
}

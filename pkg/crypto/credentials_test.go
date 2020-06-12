package crypto_test

import (
	"fmt"
	"testing"

	"github.com/pascallimeux/tools/pkg/crypto"
)

func TestGenerateHash(t *testing.T) {
	hash, err := crypto.GenerateHash("admin")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("hash = %s\n", hash)

}
func TestGenerateEncodedHash(t *testing.T) {
	encodedB64Hash, err := crypto.GenerateEncodedHash("admin")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("hash = %s\n", encodedB64Hash)

}

func TestCompareHash(t *testing.T) {
	var password = "MyPassword4Test!"
	hash, _ := crypto.GenerateHash(password)
	hashStr := string(hash)
	hash = []byte(hashStr)
	ok := crypto.CompareHashAndPassword(hash, password)
	if !ok {
		t.Error("error to verify hash password")
	}
}

func TestAesCipherWithHash(t *testing.T) {
	password := "MyPAssword2TestCipher"
	myKey1 := crypto.CreateAESKeyFromPassword(password)
	cipher1 := crypto.NewAESCipherWithKey(myKey1)
	buffer1 := []byte("This is a text to verify cipher from password!!")
	encryptBuffer, err := cipher1.Encrypt(buffer1)
	if err != nil {
		t.Error(err.Error())
	}

	myKey2 := crypto.CreateAESKeyFromPassword(password)
	cipher2 := crypto.NewAESCipherWithKey(myKey2)
	buffer2, err := cipher2.Decrypt(encryptBuffer)
	if err != nil {
		t.Error(err.Error())
	}
	if string(buffer1) != string(buffer2) {
		t.Error("error to compare 2 ciphers buffers")
	}
}

package crypto_test

import (
	"fmt"
	"testing"

	"github.com/pascallimeux/tools/pkg/crypto"
	"gopkg.in/go-playground/assert.v1"
)

func TestGenerateCipher(t *testing.T) {
	cipher, err := crypto.NewSecp256k1Cipher()
	if err != nil {
		t.Error(err.Error())
	}
	pb := cipher.GetPubKey()
	pv := cipher.GetPrivKey()
	assert.Equal(t, len(pv) == 66, true)
	assert.Equal(t, len(pb) == 132, true)
	fmt.Printf("Private Key: %s\n", pv)
	fmt.Printf("Public  key: %s\n", pb)
}

func TestGenerateCipherFromKey(t *testing.T) {
	key := "0xa6ad807cab657b88bee4e13a9c784c01e3f26a9c4f10aa5a16684afc973599de"
	cipher, err := crypto.NewSecp256k1CipherFromPrivKey(key)
	if err != nil {
		t.Error(err.Error())
	}
	pb := cipher.GetPubKey()
	pv := cipher.GetPrivKey()
	assert.Equal(t, len(pv) == 66, true)
	assert.Equal(t, len(pb) == 132, true)
	assert.Equal(t, pv == key, true)

	fmt.Printf("Private Key: %s\n", pv)
	fmt.Printf("Public  key: %s\n", pb)
}

func TestSign(t *testing.T) {
	message := "Hello world, it's a test to verify the signing mechanism..."
	cipher, err := crypto.NewSecp256k1Cipher()
	if err != nil {
		t.Error(err.Error())
	}
	sig, err := cipher.Sign(message)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, len(sig) == 130, true)
	fmt.Printf("Signature: %s\n", sig)
}

func TestCheckSign(t *testing.T) {
	message := "Hello world, it's a test to verify the signing mechanism..."
	cipher, err := crypto.NewSecp256k1Cipher()
	if err != nil {
		t.Error(err.Error())
	}
	sig, err := cipher.Sign(message)
	if err != nil {
		t.Error(err.Error())
	}
	ok, err := cipher.CheckSign(message, sig)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, ok, true)
}

func TestEncrypt(t *testing.T) {
	message := "Hello world, it's a test to verify the signing mechanism..."
	cipher, err := crypto.NewSecp256k1Cipher()
	if err != nil {
		t.Error(err.Error())
	}
	encryptedMsg, err := cipher.Encrypt(message)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("Crypted msg %s\n", encryptedMsg)
}

func TestDeCrypt(t *testing.T) {
	message := "Hello world, it's a test to verify the signing mechanism..."
	cipher, err := crypto.NewSecp256k1Cipher()
	if err != nil {
		t.Error(err.Error())
	}
	encryptedMsg, _ := cipher.Encrypt(message)
	msg, err := cipher.Decrypt(encryptedMsg)
	assert.Equal(t, message == msg, true)
}

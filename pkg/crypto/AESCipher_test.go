package crypto_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pascallimeux/tools/pkg/crypto"
)

var (
	testFolderParent = "/tmp/folder1"
	folder2          = fmt.Sprintf("%s/folder2", testFolderParent)
	testFolders      = fmt.Sprintf("%s/folder3/folder4", folder2)
	file1            = fmt.Sprintf("%s/files1.txt", testFolderParent)
	file2            = fmt.Sprintf("%s/files2.md", folder2)
)

func removeFilesAndFoldersForTests() {
	if _, err := os.Stat(testFolderParent); os.IsExist(err) {
		os.RemoveAll(testFolderParent)
	}
}

func createFilesAndFoldersForTests() error {
	if _, err := os.Stat(testFolders); os.IsNotExist(err) {
		os.MkdirAll(testFolders, os.ModePerm)
	}

	file, err := os.Create(file1)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("Hello World\n")

	file, err = os.Create(file2)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("# title1\n\n## title2\n\n\tcode\n")

	return nil
}

func Test_AEScipher(t *testing.T) {
	t.Run("Can create AES cipher", should_create_AES_cipher)
	t.Run("Can create AES cipher with key", should_create_AES_cipher_with_key)
	t.Run("Can get key", should_be_able_to_get_key)
	t.Run("Can crypt message", should_be_able_to_crypt_message)
	t.Run("Can decrypt message", should_be_able_to_decrypt_message)
}

func should_create_AES_cipher(t *testing.T) {
	cipher, err := crypto.NewAESCipher()
	if err != nil || cipher == nil {
		t.Error("Error generating AES key")
	}
}

func should_create_AES_cipher_with_key(t *testing.T) {
	cipher, err := crypto.NewAESCipher()
	key := cipher.GetEncodedKey()
	cipher2, err := crypto.NewAESCipherWithEncodedKey(key)
	if err != nil || cipher2 == nil {
		t.Error("Error generating AES key")
	}
}

func should_be_able_to_get_key(t *testing.T) {
	cipher, err := crypto.NewAESCipher()
	key := cipher.GetEncodedKey()
	if err != nil || key == "" {
		t.Error("Error getting AES key")
	}
}

func should_be_able_to_crypt_message(t *testing.T) {
	cipher, err := crypto.NewAESCipher()
	message := []byte("My name is Toto")
	crypted, err := cipher.EncryptToEncoded(message)
	if err != nil || crypted == "" {
		t.Error("Error to AES encrypte")
	}
}

func should_be_able_to_decrypt_message(t *testing.T) {
	cipher, _ := crypto.NewAESCipher()
	message := []byte("My name is Toto")
	crypted, _ := cipher.Encrypt(message)
	decrypted, _ := cipher.Decrypt(crypted)
	if string(decrypted) != string(message) {
		t.Error("Error decrypt message is different of original message")
	}
}

func TestEncryptFile(t *testing.T) {
	createFilesAndFoldersForTests()
	defer removeFilesAndFoldersForTests()
	cipher, err := crypto.NewAESCipher()
	if err != nil {
		t.Error(err.Error())
	}
	err = cipher.EncryptFile(file1)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDecryptFile(t *testing.T) {
	createFilesAndFoldersForTests()
	defer removeFilesAndFoldersForTests()
	cipher, err := crypto.NewAESCipher()
	if err != nil {
		t.Error(err.Error())
	}
	cipher.EncryptFile(file2)
	err = cipher.DecryptFile(file2)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestEncryptFolder(t *testing.T) {
	createFilesAndFoldersForTests()
	defer removeFilesAndFoldersForTests()
	cipher, err := crypto.NewAESCipher()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = cipher.EncryptFolder(folder2)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDecryptFolder(t *testing.T) {
	createFilesAndFoldersForTests()
	//defer removeFilesAndFoldersForTests()
	cipher, err := crypto.NewAESCipher()
	if err != nil {
		t.Error(err.Error())
	}
	cipher.EncryptFolder(folder2)
	_, err = cipher.DecryptFolder(folder2)
	if err != nil {
		t.Error(err.Error())
	}
}

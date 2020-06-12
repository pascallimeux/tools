package crypto_test

import (
	"os"
	"testing"
	"github.com/pascallimeux/tools/pkg/crypto"
)

var (
	privKeyPath = "/tmp/priv.pem"
	pubKeyPath  = "/tmp/pub.pem"
	message     = []byte("Shalom!")
)

func Test_RSAcipher(t *testing.T) {
	t.Run("Can create a RSA cipher in memory", should_create_RSA_cipher_in_memory)
	t.Run("Can create a RSA cipher with keys", should_create_RSA_cipher)
	t.Run("Can instanciate a RSA cipher from keys", should_instanciate_RSA_cipher)
	t.Run("can instanciate a RSA cipher with privKey", should_instanciate_RSA_cipher_with_only_privKey)
	t.Run("can instanciate a RSA cipher with pubKey", should_instanciate_RSA_cipher_with_only_pubKey)
	t.Run("Can instanciate a RSA cipher with pem memory keys", should_create_RSA_cipher_with_keys)	
	t.Run("Can encrypte a message", should_encrypt_message)
	t.Run("Can't encrypte a message", should_not_encrypt_message_less_pubKey)
	t.Run("Can decrypt a message", should_decrypt_message)
	t.Run("Can sign a message", should_sign_message)
	t.Run("Can't sign a message", should_not_sign_message_less_privKey)
	t.Run("Can check the signature of a message", should_check_message_signature)
	t.Run("Can get RSA keys", should_get_RSA_keys)
}

func should_create_RSA_cipher_in_memory(t *testing.T) {
	cipher := crypto.NewRSACipherInMemory()
	if cipher == nil {
		t.Error("Error create RSA cipher")
	}
}

func should_create_RSA_cipher(t *testing.T) {
	os.Remove("/tmp/*.pem")
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	if err != nil || cipher == nil {
		t.Error("Error create RSA cipher")
	}
}

func should_get_RSA_keys(t *testing.T) {
	cipher := crypto.NewRSACipherInMemory()
	privKey, pubKey, err := cipher.GetKeys()
	if err != nil || privKey == nil || pubKey == nil {
		t.Error("Error get RSA cipher")
	}
}

func should_create_RSA_cipher_with_keys(t *testing.T) {
	cipher := crypto.NewRSACipherInMemory()
	privKey, pubKey, err := cipher.GetKeys()
	if err != nil {
		t.Error("Error get RSA cipher")
	}
	cipher, err = crypto.NewRSACipherWithKey([]byte(privKey), []byte(pubKey))
	if err != nil || cipher == nil {
		t.Error("Error create RSA cipher")
	}
}

func should_instanciate_RSA_cipher(t *testing.T) {
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	if err != nil || cipher == nil {
		t.Error("Error create RSA cipher")
	}
}

func should_instanciate_RSA_cipher_with_only_privKey(t *testing.T) {
	pubKeyPath := ""
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	if err != nil || cipher == nil {
		t.Error("Error create RSA cipher")
	}
}

func should_instanciate_RSA_cipher_with_only_pubKey(t *testing.T) {
	privKeyPath := ""
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	if err != nil || cipher == nil {
		t.Error("Error create RSA cipher")
	}
}

func should_encrypt_message(t *testing.T) {
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	encrypted, err := cipher.EncryptToEncoded(message)
	if err != nil || encrypted == "" {
		t.Error("Error encrypte message")
	}
}

func should_not_encrypt_message_less_pubKey(t *testing.T) {
	pubKeyPath := ""
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	_, err = cipher.Encrypt(message)
	if err == nil {
		t.Error("Error encrypte message")
	}
}

func should_decrypt_message(t *testing.T) {
	cipher, _ := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	encrypted, _ := cipher.Encrypt(message)
	decrypted, _ := cipher.Decrypt(encrypted)
	if string(decrypted) != string(message) {
		t.Error("Error decrypte message")
	}
}

func should_sign_message(t *testing.T) {
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	sig, err := cipher.Sign(message)
	if err != nil || sig == nil {
		t.Error("Error sign message")
	}
}

func should_not_sign_message_less_privKey(t *testing.T) {
	privKeyPath := ""
	cipher, err := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	_, err = cipher.Sign(message)
	if err == nil {
		t.Error("Error sign message")
	}
}

func should_check_message_signature(t *testing.T) {
	cipher, _ := crypto.NewRSACipher(privKeyPath, pubKeyPath)
	sig, _ := cipher.Sign(message)
	checkSig := cipher.CheckSign(message, sig)
	if !checkSig {
		t.Error("Error check signature")
	}
}

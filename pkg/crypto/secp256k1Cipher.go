package crypto

import (
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/decred/dcrd/dcrec/secp256k1/schnorr"
)

type Secp256k1Cipher struct {
	privKey *secp256k1.PrivateKey
	pubKey  *secp256k1.PublicKey
}

//NewSecp256k1Cipher create a new Secp256k1Cipher
func NewSecp256k1Cipher() (*Secp256k1Cipher, error) {
	cipher := Secp256k1Cipher{}
	var err error
	cipher.privKey, err = secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	cipher.pubKey = cipher.privKey.PubKey()
	return &cipher, err
}

//NewSecp256k1CipherFromPrivKey create a new Secp256k1Cipher from private key in hexadecimal format
func NewSecp256k1CipherFromPrivKey(key string) (*Secp256k1Cipher, error) {
	cipher := Secp256k1Cipher{}
	pkBytes, err := hexa2Bytes(key)
	if err != nil {
		return nil, err
	}
	cipher.privKey, cipher.pubKey = secp256k1.PrivKeyFromBytes(pkBytes)
	return &cipher, err
}

// GetPubKey get public key in hexadecimal format
func (c *Secp256k1Cipher) GetPubKey() string {
	serializedPubKey := c.pubKey.SerializeUncompressed()
	return bytes2Hexa(serializedPubKey)
}

// GetPrivKey get private key in hexadecimal format
func (c *Secp256k1Cipher) GetPrivKey() string {
	serializedPrivKey := c.privKey.Serialize()
	return bytes2Hexa(serializedPrivKey)
}

//Sign sign a string message and send the sig in hexadecimal format
func (c *Secp256k1Cipher) Sign(msg string) (string, error) {
	messageHash := chainhash.HashB([]byte(msg))
	r, s, err := schnorr.Sign(c.privKey, messageHash)
	if err != nil {
		return "", err
	}
	signature := schnorr.NewSignature(r, s)
	serializedSig := signature.Serialize()
	return bytes2Hexa(serializedSig), nil
}

//CheckSign check signature in hexadecimal format
func (c *Secp256k1Cipher) CheckSign(msg, sigStr string) (bool, error) {
	sigBytes, err := hexa2Bytes(sigStr)
	if err != nil {
		return false, err
	}
	signature, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return false, err
	}
	messageHash := chainhash.HashB([]byte(msg))
	verified := schnorr.Verify(c.pubKey, messageHash, signature.R, signature.S)
	return verified, nil
}

//Encrypt encrypt a string message to hexadecimal
func (c *Secp256k1Cipher) Encrypt(msg string) (string, error) {
	encrypted, err := secp256k1.Encrypt(c.pubKey, []byte(msg))
	if err != nil {
		return "", err
	}
	return bytes2Hexa(encrypted), nil
}

//Decrypt decrypt hexadecimal message
func (c *Secp256k1Cipher) Decrypt(encryptedMsg string) (string, error) {
	encryptedBytes, err := hexa2Bytes(encryptedMsg)
	if err != nil {
		return "", err
	}
	decrypted, err := secp256k1.Decrypt(c.privKey, encryptedBytes)
	if err != nil {
		return "", nil
	}
	return string(decrypted), nil
}

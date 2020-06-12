package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

//RSACipher RSA Cipher to encrypt and decrypt
type RSACipher struct {
	privkey *rsa.PrivateKey
	pubkey  *rsa.PublicKey
}

//NewRSACipherFromPubKeyPath Instanciate new RSA cipher with public key to be used by consumer
func NewRSACipherFromPubKeyPath(pubKeyPath, privKeyPath string) (*RSACipher, error) {
	cipher := RSACipher{}
	if _, err := os.Stat(privKeyPath); err == nil {
		return nil, fmt.Errorf("Private key path already exist: %s", privKeyPath)
	}
	if _, err := os.Stat(pubKeyPath); err != nil {
		return nil, fmt.Errorf("Publmic key path does not exist: %s", pubKeyPath)
	}
	cipher.loadRSAPublicKeyFromPemFile(pubKeyPath)
	return &cipher, nil
}

//NewRSACipherInMemory Instanciate new RSA cipher in memory (without certificate files)
func NewRSACipherInMemory() *RSACipher {
	cipher := RSACipher{}
	cipher.generateRSAKeyPair()
	return &cipher
}

//NewRSACipher Instanciate new RSA cipher with private key path, public key path or nothing to be used by owner
func NewRSACipher(privKeyPath, pubKeyPath string) (*RSACipher, error) {
	cipher := RSACipher{}
	var privExist, pubExist bool
	if _, err := os.Stat(privKeyPath); err == nil {
		privExist = true
	}
	if _, err := os.Stat(pubKeyPath); err == nil {
		pubExist = true
	}
	if !privExist && !pubExist {
		cipher.generateRSAKeyPair()
		err := cipher.createRSAPrivateKeyAsPemFile(privKeyPath)
		if err != nil {
			return nil, err
		}
		err = cipher.createRSAPublicKeyAsPemFile(pubKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		if privExist {
			err := cipher.loadRSAPrivateKeyFromPemFile(privKeyPath)
			if err != nil {
				return nil, err
			}
		}
		if pubExist {
			err := cipher.loadRSAPublicKeyFromPemFile(pubKeyPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return &cipher, nil
}

//NewRSACipherWithKey Instanciate new RSA cipher with private key, public key
func NewRSACipherWithKey(privPEM, pubPEM []byte) (*RSACipher, error) {
	cipher := RSACipher{}
	if privPEM == nil && pubPEM == nil {
		return nil, fmt.Errorf("you must give in parameters at least one key (private or public)")
	}
	if privPEM != nil {
		err := cipher.loadRSAPrivateKey(privPEM)
		if err != nil {
			return nil, err
		}
	}
	if pubPEM != nil {
		err := cipher.loadRSAPublicKey(pubPEM)
		if err != nil {
			return nil, err
		}
	}
	return &cipher, nil
}

func (r *RSACipher) generateRSAKeyPair() {
	start := time.Now()
	var err error
	r.privkey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Error(fmt.Sprintf("failed generate keys: %s", err))
	}
	r.pubkey = &r.privkey.PublicKey
	log.Info(fmt.Sprintf("Generate key pair in %s", time.Since(start)))
}

func (r *RSACipher) getPrivKeyPem() []byte {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(r.privkey)
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
}

func (r *RSACipher) getPubKeyPem() ([]byte, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(r.pubkey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	), nil
}

func (r *RSACipher) createRSAPrivateKeyAsPemFile(PEMfile string) error {
	privkeyPem := r.getPrivKeyPem()
	file, err := os.Create(PEMfile)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(string(privkeyPem))
	return nil
}

func (r *RSACipher) createRSAPublicKeyAsPemFile(PEMfile string) error {
	pubkeyPem, err := r.getPubKeyPem()
	if err != nil {
		return err
	}
	file, err := os.Create(PEMfile)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(string(pubkeyPem))
	return nil
}

func (r *RSACipher) loadRSAPrivateKeyFromPemFile(PEMfile string) error {
	raw, err := ioutil.ReadFile(PEMfile)
	if err != nil {
		return err
	}
	privPEM := string(raw)
	return r.loadRSAPrivateKey([]byte(privPEM))
}

func (r *RSACipher) loadRSAPrivateKey(privPEM []byte) error {
	var err error
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return errors.New("failed to parse PEM block containing the key")
	}
	r.privkey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	return nil
}

func (r *RSACipher) loadRSAPublicKeyFromPemFile(PEMfile string) error {
	raw, err := ioutil.ReadFile(PEMfile)
	if err != nil {
		return err
	}
	pubPEM := string(raw)
	return r.loadRSAPublicKey([]byte(pubPEM))
}

func (r *RSACipher) loadRSAPublicKey(pubPEM []byte) error {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		r.pubkey = pub
		return nil
	default:
		return errors.New("Key type is not RSA")
	}
}

//GetKeys get keys in pem format
func (r *RSACipher) GetKeys() ([]byte, []byte, error) {
	privkeyPem := r.getPrivKeyPem()
	pubkeyPem, err := r.getPubKeyPem()
	return privkeyPem, pubkeyPem, err
}

func (r *RSACipher) Encrypt(message []byte) ([]byte, error) {
	start := time.Now()
	if r.pubkey == nil {
		return nil, errors.New("Can't encrypt with no public key")
	}
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, r.pubkey, message, label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("Encrypt in %s", time.Since(start)))
	return ciphertext, nil
}

//Encrypt []byte message to base 64 crypt message
func (r *RSACipher) EncryptToEncoded(message []byte) (string, error) {
	ciphertext, err := r.Encrypt(message)
	return EncodeBase64(ciphertext), err
}

func (r *RSACipher) Decrypt(ciphertext []byte) ([]byte, error) {
	start := time.Now()

	if r.privkey == nil {
		return nil, errors.New("Can't decrypt with no private key")
	}
	hash := sha256.New()
	label := []byte("")
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, r.privkey, ciphertext, label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("Decrypt in %s", time.Since(start)))
	return plainText, nil
}

//Decrypt base 64 encoded message to []byte
func (r *RSACipher) DecryptEncodedMessage(b64cipher string) ([]byte, error) {
	ciphertext, err := DecodeBase64(b64cipher)
	if err != nil {
		return nil, err
	}
	return r.Decrypt(ciphertext)
}

// Sign a []byte message and generate signature
func (r *RSACipher) Sign(message []byte) ([]byte, error) {
	if r.privkey == nil {
		return nil, errors.New("Can't sign with no private key")
	}
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(message)
	hashed := pssh.Sum(nil)
	sig, err := rsa.SignPSS(rand.Reader, r.privkey, newhash, hashed, &opts)
	return sig, err
	//b64Sign := EncodeBase64(sig)
	//return b64Sign, err
}

//CheckSign Check signature,
func (r *RSACipher) CheckSign(message []byte, sig []byte) bool {
	//signature, err := DecodeBase64(b64Sig)
	//if err != nil {
	//	return false
	//}
	var opts rsa.PSSOptions
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(message)
	hashed := pssh.Sum(nil)
	err := rsa.VerifyPSS(r.pubkey, newhash, hashed, sig, &opts)
	if err != nil {
		return false
	}
	return true
}

package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

//GenerateKeyAndCert Generate keys and certificate x509
func GenerateKeyAndCert(certPathFile, pubKeyPathFile, privKeyPathFile string) {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("GenerateKeyAndCert failed %s", err)
	}
	writeCertificate(privkey, certPathFile)
	writePrivateKey(privkey, privKeyPathFile)
	writePublicKey(privkey, pubKeyPathFile)
}

func writeCertificate(privkey *rsa.PrivateKey, certPathFile string) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"FR"},
			Organization:       []string{"ccp3"},
			OrganizationalUnit: []string{"r&d"},
			Locality:           []string{"Paris"},
			Province:           []string{"Paris"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privkey), privkey)
	if err != nil {
		log.Fatalf("certificate creation failed: %s", err)
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certfile, err := os.Create(certPathFile)
	if err != nil {
		log.Fatalf("certificate creation failed: %s", err)
	}
	defer certfile.Close()
	certfile.WriteString(out.String())
}

func writePrivateKey(privkey *rsa.PrivateKey, privKeyPathFile string) {
	out := &bytes.Buffer{}
	pem.Encode(out, pemBlockForKey(privkey))
	privKeyfile, err := os.Create(privKeyPathFile)
	if err != nil {
		log.Fatalf("private key creation failed: %s", err)
	}
	defer privKeyfile.Close()
	privKeyfile.WriteString(out.String())
}

func writePublicKey(privkey *rsa.PrivateKey, pubKeyPathFile string) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey(privkey))
	if err != nil {
		log.Fatalf("public key creation failed: %s", err)
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)
	pubKeyfile, err := os.Create(pubKeyPathFile)
	if err != nil {
		log.Fatalf("public key creation failed: %s", err)
	}
	defer pubKeyfile.Close()
	pubKeyfile.WriteString(string(pubkeyPem))
}

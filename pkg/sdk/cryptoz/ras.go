package cryptoz

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
)

func RSA() *Rsa {
	return &Rsa{}
}

type Rsa struct {
}

func (r *Rsa) Encrypt(plainTextBytes []byte, pub string) ([]byte, error) {
	pubBytes := []byte(pub)
	block, _ := pem.Decode(pubBytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey := publicKey.(*rsa.PublicKey)
	cipherTextBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainTextBytes)
	if err != nil {
		return nil, err
	}

	return cipherTextBytes, nil
}

func (r *Rsa) EncryptString(plainText string, pub string) (string, error) {
	cipherTextBytes, err := r.Encrypt([]byte(plainText), pub)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherTextBytes), nil
}

func (r *Rsa) DecryptString(cipherTextHex string, pri string) (string, error) {
	cipherTextBytes, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return "", err
	}
	plainTextBytes, err := r.Decrypt(cipherTextBytes, pri)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}

func (r *Rsa) Decrypt(cipherTextBytes []byte, pri string) ([]byte, error) {

	priBytes := []byte(pri)
	block, _ := pem.Decode(priBytes)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	plainTextBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), cipherTextBytes)
	if err != nil {
		return nil, err
	}
	return plainTextBytes, nil
}

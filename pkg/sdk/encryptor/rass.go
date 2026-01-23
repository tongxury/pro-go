package encryptor

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	PEM_BEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEM_END   = "\n-----END RSA PRIVATE KEY-----"

	PUBPEMBEGIN = "-----BEGIN PUBLIC KEY-----\n"
	PUBPEMEND   = "\n-----END PUBLIC KEY-----"
)

func RsaSign(signContent string, privateKey string, hash crypto.Hash) string {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := parsePrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature)
}

func parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	privateKey = formatPrivateKey(privateKey)
	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey.(*rsa.PrivateKey), nil
}

func formatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PEM_BEGIN) {
		privateKey = PEM_BEGIN + privateKey
	}
	if !strings.HasSuffix(privateKey, PEM_END) {
		privateKey = privateKey + PEM_END
	}
	return privateKey
}

func RsaCheckSign(signContent, sign, publicKey string, hash crypto.Hash) error {
	//hashed := sha256.Sum256([]byte(signContent))
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)

	pubKey, err := parsePublicKey(publicKey)
	if err != nil {
		return err
	}
	sig, _ := base64.RawStdEncoding.DecodeString(sign)

	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed, sig)
	if err != nil {
		return err
	}
	return nil
}

// ParsePublicKey 公钥验证
func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = formatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

func formatPublicKey(publicKey string) string {
	if !strings.HasPrefix(publicKey, PUBPEMBEGIN) {
		publicKey = PUBPEMBEGIN + publicKey
	}
	if !strings.HasSuffix(publicKey, PUBPEMEND) {
		publicKey = publicKey + PUBPEMEND
	}
	return publicKey
}

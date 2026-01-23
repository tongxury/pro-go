package encryptor

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateRsaKey(keySize int) {
	//1.生成rsa秘钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		panic(err)
	}
	//2.通过x509标准将得到的rsa私钥序列化为ASN.1的DER编码字符串
	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	//3.创建一个pem.Block结构体
	block := pem.Block{
		Type:  "rsa private key",
		Bytes: derText,
	}
	//4.通过pem将设置好的私钥数据进行编码，并写入磁盘文件
	file, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	err = pem.Encode(file, &block)
	if err != nil {
		panic(err)
	}

	// ==========公钥==================
	//1.从私钥中取出公钥
	publicKey := privateKey.PublicKey
	//2.使用x509序列化公钥为字符串
	marshalPKIXPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//3.通过公钥字符串设置到pem格式块中
	block = pem.Block{
		Type:    "rsa public key",
		Headers: nil,
		Bytes:   marshalPKIXPublicKey,
	}
	//4.pem编码
	file, err = os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	err = pem.Encode(file, &block)
	if err != nil {
		panic(err)
	}
	file.Close()
}

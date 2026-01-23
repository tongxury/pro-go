package encryptor

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
)

func DESCBC(key string) *descbc {
	return &descbc{
		key: key,
	}
}

type descbc struct {
	key string
}

func (d *descbc) Encrypt(src string) (string, error) {
	data := []byte(src)
	keyByte := []byte(d.key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return "", errors.New(err.Error())
	}
	data = PKCS5Padding(data, block.BlockSize())
	//获取CBC加密模式
	iv := keyByte //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	return fmt.Sprintf("%x", out), nil
}

func (d *descbc) Decrypt(src string) (string, error) {
	keyByte := []byte(d.key)
	data, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	iv := keyByte //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	plaintext = PKCS5UnPadding(plaintext)
	return string(plaintext), nil
}

// // CBC加密
//
//	func DESCBCEncrypt(src, key string) string {
//		data := []byte(src)
//		keyByte := []byte(key)
//		block, err := des.NewCipher(keyByte)
//		if err != nil {
//			return ""
//		}
//		data = PKCS5Padding(data, block.BlockSize())
//		//获取CBC加密模式
//		iv := keyByte //用密钥作为向量(不建议这样使用)
//		mode := cipher.NewCBCEncrypter(block, iv)
//		out := make([]byte, len(data))
//		mode.CryptBlocks(out, data)
//		return fmt.Sprintf("%X", out)
//	}
//
// // CBC解密
//func DESCBCDecrypt(src, key string) string {
//	keyByte := []byte(key)
//	data, err := hex.DecodeString(src)
//	if err != nil {
//		return ""
//	}
//	block, err := des.NewCipher(keyByte)
//	if err != nil {
//		return ""
//	}
//	iv := keyByte //用密钥作为向量(不建议这样使用)
//	mode := cipher.NewCBCDecrypter(block, iv)
//	plaintext := make([]byte, len(data))
//	mode.CryptBlocks(plaintext, data)
//	plaintext = PKCS5UnPadding(plaintext)
//	return string(plaintext)
//}

// ECB加密 key固定8位
func DESECBEncrypt(src, key string) string {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return ""
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return fmt.Sprintf("%X", out)
}

// // ECB解密
func DESECBDecrypt(src, key string) string {
	data, err := hex.DecodeString(src)
	if err != nil {
		return ""
	}
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return ""
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return string(out)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("unpadding error")
	}
	return origData[:(length - unpadding)]
}

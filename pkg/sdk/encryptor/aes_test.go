package encryptor

import (
	"fmt"
	"store/pkg/sdk/conv"
	"testing"
)

func TestName(t *testing.T) {
	//key := fmt.Sprintf("%x", sha256.Sum256([]byte("382450801@qq.com")))[:8]
	//key := fmt.Sprintf("%x", sha256.Sum256([]byte("ec14f21747048b63")))[:16]

	//c7400d4ef220789ea3adf6f91b69ae57
	//a := ECBEncrypt([]byte("Aa1774566!"), []byte("f769b57395e9d73a"))
	a := ECBEncrypt(conv.M2B(map[string]string{
		"email": "2210515001@email.szu.edu.cn", "password": "Aa1774566!",
	}), []byte("43cbadc3d65b4a6e"))

	//d4e293ad418c5d2d4b3fe005c27cf3d1

	fmt.Println(fmt.Sprintf("%x", a))
}

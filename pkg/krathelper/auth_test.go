package krathelper

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	//fmt.Println(GenerateToken(5068))
	//fmt.Println(GenerateToken(1665368111286104066))
	//fmt.Println(GenerateToken(1684137060545933313))
	//fmt.Println(GenerateTokenV2("690d6a669e5c05462c0e4165"))
	//fmt.Println(GenerateTokenV2("10002232"))
	fmt.Println(GenerateTokenV2("69362f973aa0830f527772f2"))
	fmt.Println(GenerateTokenV2("693b91f82b271bf02ac1e624"))

	fmt.Println(GenerateTokenV2("6980477237ca92ff5dbf7cd7"))
	fmt.Println(GenerateTokenV2("698722f437ca92ff5dbf7cd8"))
	//claims, _ := ParseClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0cyI6MTc2NDUyNTQ1MSwidXNlcl9pZCI6IiJ9.Dx687uOFHKX8_JxGj08I1oRauNGpQ1oRPEyb5Q1_Go4", SecretSignKey)

	//fmt.Println(claims)
}

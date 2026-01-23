package conv

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func HexToInt64(hex string) (int64, error) {
	val := strings.ReplaceAll(hex, "0x", "")
	return strconv.ParseInt(val, 16, 64)
}

func TrimLeftZeros(val string) string {

	if strings.HasPrefix(val, "0x") {
		val = val[2:]
	}

	var trimFrom int
	for i, v := range val {
		if v != '0' {
			trimFrom = i
			break
		}
	}

	return "0x" + val[trimFrom:]
}

func HexToBigInt(hex string) *big.Int {

	if strings.HasPrefix(hex, "0x") {
		hex = hex[2:]
	}

	n := new(big.Int)
	if hex == "" {
		return n
	}

	n, _ = n.SetString(hex, 16)

	return n
}

func MustHexToBytes(hexStr string) []byte {

	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}

	decodeBytes, _ := hex.DecodeString(hexStr)

	return decodeBytes
}

func BytesToHex(src []byte) string {
	return fmt.Sprintf("%x", src)
}

//func HexToBytes32(hexStr string) ([32]byte, error) {
//
//	rsp := [32]byte{}
//
//	bytes, err := HexToBytes(hexStr)
//	if err != nil {
//		return rsp, err
//	}
//
//	copy(rsp[:], bytes)
//
//	return rsp, nil
//}

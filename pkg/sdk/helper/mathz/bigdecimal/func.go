package bigdecimal

import (
	"math"
	"math/big"
	"strings"
)

func (t BigDecimal) decimalPow10() *big.Float {
	return new(big.Float).SetInt64(int64(math.Pow10(t.decimals)))
}

func (t BigDecimal) Float64() float64 {
	z := new(big.Float).Quo(t.raw, t.decimalPow10())
	v, _ := z.Float64()
	return v
}

func (t BigDecimal) Int64() int64 {
	z := new(big.Float).Quo(t.raw, t.decimalPow10())
	v, _ := z.Int64()
	return v
}

// 去除末尾无效的零，与小数点
func (t BigDecimal) trimTrailingZeros(s string) string {
	if strings.Contains(s, ".") {
		// 拆分整数部分和小数部分
		parts := strings.Split(s, ".")
		integral := parts[0]
		decimal := parts[1]

		// 去除小数部分末尾的零，保留合法的数位
		decimal = strings.TrimRight(decimal, "0")

		// 如果小数部分去掉零后为空，返回仅有整数部分的字符串
		if decimal == "" {
			return integral
		}

		// 返回整合后的字符串
		return integral + "." + decimal
	}
	return s // 如果没有小数点，直接返回原字符串
}

func (t BigDecimal) Raw() string {

	v, _ := t.raw.Int(nil)

	return v.String()
	//return t.trimTrailingZeros(t.raw.Text('f', 20))
}

//func (t BigDecimal) RawInt() string {
//
//	v, _ := t.raw.Int(nil)
//
//	return v.String()
//}

func (t BigDecimal) GreaterThenZero() bool {
	v, _ := t.raw.Float64()
	return v > 0
}

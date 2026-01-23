package bigdecimal

import (
	"math"
	"math/big"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/mathz"
)

func New[T mathz.Number, V mathz.Int](value T, decimals V) BigDecimal {

	val := new(big.Float).SetFloat64(float64(value))
	dec := big.NewFloat(math.Pow10(int(decimals)))

	return BigDecimal{
		raw:      new(big.Float).Mul(val, dec),
		decimals: int(decimals),
	}
}

func FromOrigin[T any, V mathz.Int](value T, decimals V) BigDecimal {

	val, _ := new(big.Float).SetString(conv.Str(value))

	return BigDecimal{
		raw:      val,
		decimals: int(decimals),
	}
}

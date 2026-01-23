package decimal

import (
	"store/pkg/sdk/helper/mathz"
)

func New[T mathz.Number, V mathz.Int](value T, decimals V) Decimal {

	// val := new(big.Float).SetFloat64(float64(value))
	// dec := big.NewFloat(math.Pow10(int(decimals)))

	return Decimal{
		// Raw:      new(big.Float).Mul(val, dec),
		// Decimals: int(decimals),
	}
}

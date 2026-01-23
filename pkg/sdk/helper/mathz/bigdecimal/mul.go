package bigdecimal

import "math/big"

func (t BigDecimal) MulFloat64(v float64) BigDecimal {

	return BigDecimal{
		raw:      new(big.Float).Mul(t.raw, new(big.Float).SetFloat64(v)),
		decimals: t.decimals,
	}
}

package bigdecimal

import "math/big"

func (t BigDecimal) Sub(v BigDecimal) BigDecimal {

	if t.decimals != v.decimals {
		panic("decimals do not match")
	}

	return BigDecimal{
		raw:      new(big.Float).Sub(t.raw, v.raw),
		decimals: t.decimals,
	}
}

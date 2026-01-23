package bigdecimal

import (
	"math/big"
)

type BigDecimal struct {
	raw      *big.Float
	decimals int
}

//func (t BigDecimal) MarshalJSON() ([]byte, error) {
//	return t.Float64()
//}

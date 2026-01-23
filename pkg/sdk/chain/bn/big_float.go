package bn

import (
	"math"
	"math/big"
)

type BigFloat struct {
	raw      *big.Float
	decimals int
}

func MustFloat(bigNumber any, decimals int) *BigFloat {
	bn, _ := NewFloat(bigNumber, decimals)
	return bn
}

func MustFloatFromFloat(number float64, decimals int) *BigFloat {

	bigNumber := number * math.Pow10(decimals)

	bn, _ := NewFloat(bigNumber, decimals)
	return bn
}

func NewFloat(bigNumber any, decimals int) (*BigFloat, bool) {

	var bn *big.Float
	var ok bool

	switch t := bigNumber.(type) {
	case string:
		bn, ok = new(big.Float).SetString(t)
	case float64:
		bn = new(big.Float).SetFloat64(t)
		ok = true
	case uint64:
		bn = new(big.Float).SetInt64(int64(t))
		ok = true
	case int64:
		bn = new(big.Float).SetInt64(t)
		ok = true
	case int8:
		bn = new(big.Float).SetInt64(int64(t))
		ok = true
	case int:
		bn = new(big.Float).SetInt64(int64(t))
		ok = true
	case *big.Float:
		bn = new(big.Float).Set(t)
		ok = true
	case big.Float:
		bn = new(big.Float).Set(&t)
		ok = true
	}

	if !ok {
		return nil, false
	}

	return &BigFloat{raw: bn, decimals: decimals}, true
}

func (t *BigFloat) BigVal(prec ...int) string {

	precVal := 0
	if len(prec) > 0 {
		precVal = prec[0]
	}

	return t.raw.Text('f', precVal)
}

func (t *BigFloat) Div(value *BigFloat) float64 {
	return t.Val() / value.Val()
}

func (t *BigFloat) Val() float64 {

	d := new(big.Float).SetFloat64(math.Pow10(t.decimals))
	z := new(big.Float).Quo(t.raw, d)

	v, _ := z.Float64()
	return v
}

func (t *BigFloat) Int64Val() int64 {

	d := new(big.Float).SetFloat64(math.Pow10(t.decimals))
	z := new(big.Float).Quo(t.raw, d)

	v, _ := z.Int64()
	return v
}

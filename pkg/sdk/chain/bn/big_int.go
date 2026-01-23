package bn

import (
	"math"
	"math/big"
	"store/pkg/sdk/conv"
)

type BigInt struct {
	raw      *big.Int
	decimals int
}

func MustInt(bigNumber any, decimals ...int) *BigInt {
	bn, _ := NewInt(bigNumber, decimals...)

	return bn
}

func MustIntFromFloat(number float64, decimals int) *BigInt {

	val, _ := new(big.Float).Mul(big.NewFloat(number), big.NewFloat(math.Pow10(decimals))).Int64()

	return MustInt(val, decimals)
}

func MustIntFromInt(number int64, decimals ...int) *BigInt {

	dec := 0
	if len(decimals) > 0 {
		dec = decimals[0]
	}

	bigNumber := conv.Int64(float64(number) * math.Pow10(dec))

	bn, _ := NewInt(bigNumber, decimals...)
	return bn
}

func NewInt(bigNumber any, decimals ...int) (*BigInt, bool) {

	dec := 0
	if len(decimals) > 0 {
		dec = decimals[0]
	}

	var bn *big.Int
	var ok bool

	switch t := bigNumber.(type) {
	case string:
		bn, ok = new(big.Int).SetString(t, 10)
	case uint64:
		bn = new(big.Int).SetInt64(int64(t))
		ok = true
	case int64:
		bn = new(big.Int).SetInt64(t)
		ok = true
	case int8:
		bn = new(big.Int).SetInt64(int64(t))
		ok = true
	case int:
		bn = new(big.Int).SetInt64(int64(t))
		ok = true
	case *big.Int:
		bn = new(big.Int).Set(t)
		ok = true
	case big.Int:
		bn = new(big.Int).Set(&t)
		ok = true
	}

	if !ok {
		return nil, false
	}

	return &BigInt{raw: bn, decimals: dec}, true
}

func (t *BigInt) RawStr() string {
	if t == nil {
		return "0"
	}
	return t.raw.Text(10)
}

func (t *BigInt) RawVal() int64 {
	if t == nil {
		return 0
	}
	return t.raw.Int64()
}

func (t *BigInt) Int64() int64 {
	if t == nil {
		return 0
	}
	d := new(big.Float).SetInt64(int64(math.Pow10(t.decimals)))
	z := new(big.Float).Quo(new(big.Float).SetInt(t.raw), d)

	v, _ := z.Int64()
	return v
}

func (t *BigInt) Div(value *BigInt) float64 {
	if t == nil {
		return 0
	}
	return float64(t.Val()) / float64(value.Val())
}

//func (t *BigInt) Mul(value float64) *BigInt {
//	if t == nil {
//		return nil
//	}
//
//	return &BigInt{
//		raw:      new(big.Float).Mul(new(big.Float).SetInt(t.raw), big.NewFloat(value)).Int(),
//		decimals: t.decimals,
//	}
//}

func (t *BigInt) Val() float64 {

	if t == nil {
		return 0
	}

	d := new(big.Float).SetInt64(int64(math.Pow10(t.decimals)))
	z := new(big.Float).Quo(new(big.Float).SetInt(t.raw), d)

	v, _ := z.Float64()
	return v
}

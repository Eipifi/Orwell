package foo
import (
    "bytes"
    "math/big"
)

func Equal(a U256, b U256) bool {
    return Compare(a, b) == 0
}

func Compare(a, b U256) int {
    return bytes.Compare(a[:], b[:])
}

func BigToU256(b *big.Int) (r U256) {
    r.SetBytes(b.Bytes())
    return
}

func (i *U256) Big() *big.Int {
    return new(big.Int).SetBytes(i[:])
}

func (i *U256) MulDiv64(numerator, denominator uint64) U256 {
    d := i.Big()
    d = d.Mul(d, new(big.Int).SetUint64(numerator))
    d = d.Div(d, new(big.Int).SetUint64(denominator))
    return BigToU256(d)
}

func (i *U256) Invert256() U256 {
    n := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
    n = n.Div(n, i.Big())
    return BigToU256(n)
}

func (i *U256) Plus(v U256) U256 {
    d := i.Big()
    d = d.Add(d, v.Big())
    return BigToU256(d)
}

func (i *U256) Minus(v U256) U256 {
    d := i.Big()
    d = d.Sub(d, v.Big())
    return BigToU256(d)
}
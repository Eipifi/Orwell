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

func (i *U256) SetBig(b *big.Int) {
    i.SetBytes(b.Bytes())
}

func (i *U256) Big() *big.Int {
    return new(big.Int).SetBytes(i[:])
}

func (i *U256) MulDiv64(numerator, denominator uint64) {
    d := i.Big()
    d = d.Mul(d, new(big.Int).SetUint64(numerator))
    d = d.Div(d, new(big.Int).SetUint64(denominator))
    i.SetBig(d)
}

func (i *U256) Invert256() {
    n := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
    n = n.Div(n, i.Big())
    i.SetBig(n)
}

func (i *U256) Add(v U256) {
    d := i.Big()
    d = d.Add(d, v.Big())
    i.SetBig(d)
}

func (i *U256) Sub(v U256) {
    d := i.Big()
    d = d.Sub(d, v.Big())
    i.SetBig(d)
}
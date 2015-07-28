package foo
import (
    "io"
    "orwell/lib/butils"
    "encoding/hex"
)

const U256_BITS = 256
const U256_BYTES = U256_BITS / 8

type U256 [U256_BYTES]byte

var ZERO = U256{}
var ONE, _ = FromHex("0000000000000000000000000000000000000000000000000000000000000001")
var MAX, _ = FromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

func (i *U256) Read(r io.Reader) error {
    return butils.ReadFull(r, i[:])
}

func (i *U256) Write(w io.Writer) error {
    return butils.WriteFull(w, i[:])
}

func (i U256) String() string {
    return hex.EncodeToString(i[:])
}
package butils
import (
    "io"
    "encoding/hex"
    "bytes"
    "fmt"
    "strconv"
    "errors"
)

const (
    UINT256_LENGTH_BYTES = 32
    UINT256_LENGTH_BITS = UINT256_LENGTH_BYTES * 8
)

type Uint256 [UINT256_LENGTH_BYTES]byte

func (i *Uint256) Read(r io.Reader) error {
    return ReadFull(r, i[:])
}

func (i *Uint256) Write(w io.Writer) error {
    return WriteFull(w, i[:])
}

func Equal(a Uint256, b Uint256) bool {
    return Compare(a, b) == 0
}

func Compare(a, b Uint256) int {
    return bytes.Compare(a[:], b[:])
}

func (i Uint256) String() string {
    return hex.EncodeToString(i[:])
}

func (i Uint256) Bin() string {
    str := ""
    for x := 0; x < UINT256_LENGTH_BYTES; x += 1 {
        str += fmt.Sprintf("%08s", strconv.FormatUint(uint64(i[x]), 2))
    }
    return str
}

func FromHex(s string) (i Uint256, err error) {
    if len(s) != 64 { return i, errors.New("invalid hex length") }
    d, err := hex.DecodeString(s)
    if err != nil { return }
    copy(i[:], d)
    return
}

func (i *Uint256) SetBit(bit uint8) {
    s := bit / 8
    r := bit % 8
    i[s] |= 0x80 >> r
}

func (i *Uint256) ShiftR(n uint8) {
    for ; n > 0; n -= 1 {
        for x := UINT256_LENGTH_BYTES - 1; x > 0; x -= 1 {
            i[x] = (i[x] >> 1) | ((i[x-1] & 0x01) << 7)
        }
        i[0] = i[0] >> 1
    }
}
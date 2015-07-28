package foo
import "encoding/hex"

func (i *U256) SetBytes(value []byte) {
    if len(value) > U256_BYTES {
        copy(i[:], value[len(value)-U256_BYTES:])
    } else {
        copy(i[U256_BYTES-len(value):], value)
    }
}

func FromHex(value string) (result U256, err error) {
    buf, err := hex.DecodeString(value)
    if err != nil { return }
    result.SetBytes(buf)
    return
}


package foo
import "encoding/hex"

func (i *U256) SetBytes(value []byte) {
    if len(value) > U256_BYTES {
        copy(i[:], value[len(value)-U256_BYTES:])
    } else {
        copy(i[U256_BYTES-len(value):], value)
    }
}

func FromHex(value string) (result U256) {
    buf, err := hex.DecodeString(value)
    if err != nil {
        panic(err)
    }
    result.SetBytes(buf)
    return
}


package protocol
import (
    "orwell/orlib/sig"
    "math/rand"
)

type Address struct {
    IP [16]byte
    Port uint16
    Nonce uint64
}

func (a *Address) Id() *sig.ID {
    var b [24]byte
    copy(b[0:16], a.IP[:])
    ByteOrder.PutUint64(b[16:24], a.Nonce)
    return sig.Hash(b[:])
}

type TTL uint8
type Token uint64

func RandomToken() Token {
    v := rand.Uint32() << 4 + rand.Uint32()
    return Token(v)
}
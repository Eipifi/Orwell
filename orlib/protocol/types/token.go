package types
import (
    "orwell/orlib/comm"
    "math/rand"
)

type Token uint64

func (t *Token) Read(r *comm.Reader) (err error) {
    var v uint64
    v, err = r.ReadUint64()
    *t = Token(v)
    return
}

func (t *Token) Write(w *comm.Writer) {
    w.WriteUint64(uint64(*t))
}

func RandomToken() Token {
    v := rand.Uint32() << 4 + rand.Uint32()
    return Token(v)
}
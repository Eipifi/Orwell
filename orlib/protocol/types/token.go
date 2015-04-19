package types
import (
    "orwell/orlib/comm"
    "crypto/rand"
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

var randomSrc *comm.Reader = comm.NewReader(rand.Reader)

func RandomToken() Token {
    v, _ := randomSrc.ReadUint64()
    return Token(v)
}
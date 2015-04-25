package common
import (
    "io"
    "orwell/orlib2/butils"
    "crypto/rand"
)

type Token uint64

func (t *Token) Read(r io.Reader) error {
    val, err := butils.ReadUint64(r)
    if err != nil { return err }
    *t = Token(val)
    return
}

func (t *Token) Write(w io.Writer) error {
    return butils.WriteUint64(w, uint64(*t))
}

func NewRandomToken() Token {
    v, _ := butils.ReadUint64(rand.Reader)
    return Token(v)
}
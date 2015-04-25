package orcache
import (
    "orwell/orlib2/protocol/common"
    "io"
)

type CardNotFound struct {
    Token common.Token
    TTL common.TTL
}

func (m *CardNotFound) Read(r io.Reader) (err error) {
    if err = m.Token.Read(r); err != nil { return }
    if err = m.TTL.Read(r); err != nil { return }
    return
}

func (m *CardNotFound) Write(w io.Writer) (err error) {
    if err = m.Token.Write(w); err != nil { return }
    return m.TTL.Write(w)
}
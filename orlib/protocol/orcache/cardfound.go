package orcache
import (
    "orwell/orlib/protocol/types"
    "orwell/orlib/comm"
)

type CardFound struct {
    Token types.Token
    Card []byte
}

func (m *CardFound) Read(r *comm.Reader) (err error) {
    if err = m.Token.Read(r); err != nil { return }
    if m.Card, err = r.ReadVarBytes(); err != nil { return }
    return
}

func (m *CardFound) Write(w *comm.Writer) {
    m.Token.Write(w)
    w.WriteVarBytes(m.Card)
}
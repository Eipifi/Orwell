package orcache
import (
    "orwell/orlib/comm"
    "orwell/orlib/protocol/types"
)

type CardNotFound struct {
    Token types.Token
    TTL types.TTL
}

func (m *CardNotFound) Read(r *comm.Reader) (err error) {
    if err = m.Token.Read(r); err != nil { return }
    if err = m.TTL.Read(r); err != nil { return }
    return
}

func (m *CardNotFound) Write(w *comm.Writer) {
    m.Token.Write(w)
    m.TTL.Write(w)
}

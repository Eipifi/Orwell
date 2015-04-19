package orcache
import (
    "orwell/orlib/protocol/types"
    "orwell/orlib/comm"
)

type Publish struct {
    Token types.Token
    TTL types.TTL
    Card []byte
}

func (p *Publish) Read(r *comm.Reader) (err error) {
    if err = p.Token.Read(r); err != nil { return }
    if err = p.TTL.Read(r); err != nil { return }
    if p.Card, err = r.ReadVarBytes(); err != nil { return }
    return
}

func (p *Publish) Write(w *comm.Writer) {
    p.Token.Write(w)
    p.TTL.Write(w)
    w.WriteVarBytes(p.Card)
}
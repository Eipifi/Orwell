package orcache
import (
    "orwell/orlib/protocol/types"
    "orwell/orlib/comm"
)

type Published struct {
    Token types.Token
    TTL types.TTL
}

func (p *Published) Read(r *comm.Reader) (err error) {
    if err = p.Token.Read(r); err != nil { return }
    if err = p.TTL.Read(r); err != nil { return }
    return
}

func (p *Published) Write(w *comm.Writer) {
    p.Token.Write(w)
    p.TTL.Write(w)
}

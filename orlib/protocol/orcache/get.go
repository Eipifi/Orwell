package orcache
import (
    "orwell/orlib/protocol/types"
    "orwell/orlib/comm"
)

type Get struct {
    Token types.Token
    TTL types.TTL
    ID *types.ID
    Version uint64
}

func (m *Get) Read(r *comm.Reader) (err error) {
    m.ID = &types.ID{}
    if err = m.Token.Read(r); err != nil { return }
    if err = m.TTL.Read(r); err != nil { return }
    if err = m.ID.Read(r); err != nil { return }
    if m.Version, err = r.ReadVaruint(); err != nil { return }
    return
}

func (m *Get) Write(w *comm.Writer) {
    m.Token.Write(w)
    m.TTL.Write(w)
    m.ID.Write(w)
    w.WriteVaruint(m.Version)
}
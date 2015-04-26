package orcache
import (
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
    "io"
    "orwell/orlib/butils"
)

type Get struct {
    Token common.Token
    TTL common.TTL
    ID *hash.ID
    Version uint64
}

func (g *Get) Read(r io.Reader) (err error) {
    if err = g.Token.Read(r); err != nil { return }
    if err = g.TTL.Read(r); err != nil { return }
    g.ID = &hash.ID{}
    if err = g.ID.Read(r); err != nil { return }
    g.Version, err = butils.ReadVarUint(r)
    return
}

func (g *Get) Write(w io.Writer) (err error) {
    if err = g.Token.Write(w); err != nil { return }
    if err = g.TTL.Write(w); err != nil { return }
    if err = g.ID.Write(w); err != nil { return }
    return butils.WriteVarUint(w, g.Version)
}
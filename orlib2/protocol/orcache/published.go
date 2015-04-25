package orcache
import (
    "orwell/orlib2/protocol/common"
    "io"
)

type Published struct {
    Token common.Token
    TTL common.TTL
}

func (p *Published) Read(r io.Reader) (err error) {
    if err = p.Token.Read(r); err != nil { return }
    if err = p.TTL.Read(r); err != nil { return }
    return
}

func (p *Published) Write(w io.Writer) (err error) {
    if err = p.Token.Write(w); err != nil { return }
    return p.TTL.Write(w)
}
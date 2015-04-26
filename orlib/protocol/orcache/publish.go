package orcache
import (
    "orwell/orlib/protocol/common"
    "io"
    "orwell/orlib/crypto/card"
    "orwell/orlib/butils"
)

type Publish struct {
    Token common.Token
    TTL common.TTL
    Card *card.Card
}

func (p *Publish) Read(r io.Reader) (err error) {
    if err = p.Token.Read(r); err != nil { return }
    if err = p.TTL.Read(r); err != nil { return }
    var c []byte
    if c, err = butils.ReadVarBytes(r); err != nil { return }
    p.Card = &card.Card{}
    return p.Card.ReadBytes(c)
}

func (p *Publish) Write(w io.Writer) (err error) {
    if err = p.Token.Write(w); err != nil { return }
    if err = p.TTL.Write(w); err != nil { return }
    var card []byte
    if card, err = p.Card.WriteBytes(); err != nil { return }
    return butils.WriteVarBytes(w, card)
}
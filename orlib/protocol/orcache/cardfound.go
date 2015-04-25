package orcache
import (
    "io"
    "orwell/orlib/protocol/common"
    "orwell/orlib/crypto/card"
    "orwell/orlib/butils"
)

type CardFound struct {
    Token common.Token
    Card *card.Card
}

func (m *CardFound) Read(r io.Reader) (err error) {
    if err = m.Token.Read(r); err != nil { return }
    var card []byte
    if card, err = butils.ReadVarBytes(r); err != nil { return }
    m.Card = &card.Card{}
    return m.Card.ReadBytes(card)
}

func (m *CardFound) Write(w io.Writer) error {
    m.Token.Write(w)
    card, err := m.Card.WriteBytes()
    if err != nil { return err }
    return butils.WriteVarBytes(w, card)
}
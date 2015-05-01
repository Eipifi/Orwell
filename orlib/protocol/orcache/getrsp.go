package orcache
import (
    "io"
    "orwell/orlib/protocol/common"
    "orwell/orlib/crypto/card"
    "orwell/orlib/butils"
)

type GetRsp struct {
    Token common.Token
    TTL common.TTL
    Card *card.Card
}

func (m *GetRsp) Read(r io.Reader) (err error) {
    if err = m.Token.Read(r); err != nil { return }
    if err = m.TTL.Read(r); err != nil { return }
    var buf []byte
    buf, err = butils.ReadVarBytes(r)
    if err != nil { return }
    if len(buf) > 0 {
        m.Card = &card.Card{}
        return m.Card.ReadBytes(buf)
    } else {
        m.Card = nil
        return
    }
}

func (m *GetRsp) Write(w io.Writer) (err error) {
    if err = m.Token.Write(w); err != nil { return }
    if err = m.TTL.Write(w); err != nil { return }
    if m.Card == nil {
        tmp := make([]byte, 0)
        return butils.WriteVarBytes(w, tmp)
    } else {
        card, err := m.Card.WriteBytes()
        if err != nil { return err }
        return butils.WriteVarBytes(w, card)
    }
}

func (m *GetRsp) GetToken() common.Token {
    return m.Token
}
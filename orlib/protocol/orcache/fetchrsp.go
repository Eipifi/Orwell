package orcache
import (
    "io"
    "orwell/orlib/protocol/common"
    "orwell/orlib/crypto/card"
    "orwell/orlib/butils"
)

type FetchRsp struct {
    Token common.Token
    TTL common.TTL
    Card *card.Card
}

func (f *FetchRsp) Read(r io.Reader) (err error) {
    if err = f.Token.Read(r); err != nil { return }
    if err = f.TTL.Read(r); err != nil { return }
    var buf []byte
    buf, err = butils.ReadVarBytes(r)
    if err != nil { return }
    if len(buf) > 0 {
        f.Card = &card.Card{}
        return f.Card.ReadBytes(buf)
    } else {
        f.Card = nil
        return
    }
}

func (f *FetchRsp) Write(w io.Writer) (err error) {
    if err = f.Token.Write(w); err != nil { return }
    if err = f.TTL.Write(w); err != nil { return }
    if f.Card == nil {
        tmp := make([]byte, 0)
        return butils.WriteVarBytes(w, tmp)
    } else {
        card, err := f.Card.WriteBytes()
        if err != nil { return err }
        return butils.WriteVarBytes(w, card)
    }
}

func (f *FetchRsp) GetToken() common.Token {
    return f.Token
}
package orcache
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
)

type GetReq struct {
    Token uint64
    TTL uint8
    Key foo.U256
}

func (m *GetReq) Read(r io.Reader) (err error) {
    if m.Token, err = butils.ReadUint64(r); err != nil { return }
    if m.TTL, err = butils.ReadUint8(r); err != nil { return }
    if err = m.Key.Read(r); err != nil { return }
    return
}

func (m *GetReq) Write(w io.Writer) (err error) {
    if err = butils.WriteUint64(w, m.Token); err != nil { return }
    if err = butils.WriteUint8(w, m.TTL); err != nil { return }
    if err = m.Key.Write(w); err != nil { return }
    return
}

type GetResp struct {
    TTL uint8
    Card *Card
}

func (m *GetResp) Read(r io.Reader) (err error) {
    if m.TTL, err = butils.ReadUint8(r); err != nil { return }
    c := &Card{}
    flag, err := butils.ReadOptional(r, c)
    if err != nil { return }
    if flag != 0x00 { m.Card = c }
    return
}

func (m *GetResp) Write(w io.Writer) (err error) {
    if err = butils.WriteUint8(w, m.TTL); err != nil { return }
    if err = butils.WriteOptional(w, m.Card); err != nil { return }
    return
}
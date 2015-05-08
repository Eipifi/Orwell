package orcache
import (
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
    "io"
    "orwell/orlib/butils"
)

type FetchReq struct {
    Token common.Token
    TTL common.TTL
    ID hash.ID
    Version uint64
}

func (*FetchReq) Code() byte { return 0x02 }

func (f *FetchReq) Read(r io.Reader) (err error) {
    if err = f.Token.Read(r); err != nil { return }
    if err = f.TTL.Read(r); err != nil { return }
    if err = f.ID.Read(r); err != nil { return }
    f.Version, err = butils.ReadVarUint(r)
    return
}

func (f *FetchReq) Write(w io.Writer) (err error) {
    if err = f.Token.Write(w); err != nil { return }
    if err = f.TTL.Write(w); err != nil { return }
    if err = f.ID.Write(w); err != nil { return }
    return butils.WriteVarUint(w, f.Version)
}

func (f *FetchReq) GetToken() common.Token {
    return f.Token
}
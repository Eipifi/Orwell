package orcache
import (
    "orwell/orlib/protocol/common"
    "io"
)

type PublishRsp struct {
    Token common.Token
    TTL common.TTL
}

func (*PublishRsp) Code() byte { return 0x83 }

func (p *PublishRsp) Read(r io.Reader) (err error) {
    if err = p.Token.Read(r); err != nil { return }
    if err = p.TTL.Read(r); err != nil { return }
    return
}

func (p *PublishRsp) Write(w io.Writer) (err error) {
    if err = p.Token.Write(w); err != nil { return }
    return p.TTL.Write(w)
}

func (p *PublishRsp) GetToken() common.Token {
    return p.Token
}
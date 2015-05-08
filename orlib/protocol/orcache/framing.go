package orcache
import (
    "orwell/orlib/butils"
    "io"
    "errors"
    "orwell/orlib/protocol/common"
)

const Magic uint32 = 0xcafebabe
const Version uint64 = 1

var ErrUnexpectedCode = errors.New("Unexpected message code")

type Message interface {
    butils.Chunk
    Code() byte
}

type TokenMessage interface {
    Message
    GetToken() common.Token
}

func ReadMessage(r io.Reader, msg Message) error {
    code, err := butils.ReadByte(r)
    if err != nil { return err }
    if code != msg.Code() { return ErrUnexpectedCode }
    return msg.Read(r)
}

func WriteMessage(w io.Writer, msg Message) (err error) {
    if err = butils.WriteByte(w, msg.Code()); err != nil { return }
    return msg.Write(w)
}

func ReadAnyMessage(r io.Reader) (msg Message, err error) {
    code, err := butils.ReadByte(r)
    if err != nil { return }
    switch code {
        case (*Handshake)(nil).Code():      msg = &Handshake{}
        case (*HandshakeAck)(nil).Code():   msg = &HandshakeAck{}
        case (*FetchReq)(nil).Code():       msg = &FetchReq{}
        case (*FetchRsp)(nil).Code():       msg = &FetchRsp{}
        case (*PublishReq)(nil).Code():     msg = &PublishReq{}
        case (*PublishRsp)(nil).Code():     msg = &PublishRsp{}
        case (*PeersReq)(nil).Code():       msg = &PeersReq{}
        case (*PeersRsp)(nil).Code():       msg = &PeersRsp{}
        default: return nil, ErrUnexpectedCode
    }
    return msg, msg.Read(r)
}
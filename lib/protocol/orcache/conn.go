package orcache
import (
    "net"
    "errors"
    "bytes"
    "orwell/lib/obp"
    "orwell/lib/butils"
    "fmt"
)

func Encode(msg obp.Msg) ([]byte, error) {
    if code, ok := MsgToCode(msg); ok {
        buf := &bytes.Buffer{}
        buf.WriteByte(code)
        err := msg.Write(buf)
        return buf.Bytes(), err
    } else { return nil, errors.New("Unknown message type") }
}

func Decode(data []byte) (obp.Msg, error) {
    if len(data) == 0 { return nil, errors.New("Cannot decode an empty byte slice into a message") }
    if msg, ok := CodeToMsg(data[0]); ok {
        err := butils.ReadAllInto(msg, data[1:])
        return msg, err
    }
    return nil, errors.New(fmt.Sprintf("Unknown message code: %v", data[0]))
}

func Connection(socket net.Conn) *obp.MsgConn {
    return obp.NewMsgConn(obp.New(socket), Encode, Decode)
}

func CodeToMsg(code byte) (obp.Msg, bool) {
    if code == MSG_CODE_GET_REQ { return &GetReq{}, true }
    if code == MSG_CODE_GET_RSP { return &GetResp{}, true }
    return nil, false
}

func MsgToCode(msg obp.Msg) (byte, bool) {
    switch msg.(type) {
        case (*GetReq): return MSG_CODE_GET_REQ, true
        case (*GetResp): return MSG_CODE_GET_RSP, true
    }
    return 0x00, false
}

const MSG_CODE_GET_REQ byte = 0x01
const MSG_CODE_GET_RSP byte = 0x81
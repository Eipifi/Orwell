package orchain
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
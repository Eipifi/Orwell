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

type HsVerifier func(*HandshakeReq) error

func Connect(socket net.Conn, hs_local *HandshakeReq, v HsVerifier) (conn *obp.MsgConn, err error) {

    // Create a message-based connection from the raw socket
    conn = obp.NewMsgConn(obp.New(socket), Encode, Decode)

    // Create the handshake handler function
    // This function will be called when a remote sends a request (we expect a HandshakeReq).
    handler := func(msg obp.Msg) (obp.Msg, error) {
        hs, ok := msg.(*HandshakeReq)
        if ! ok { return nil, errors.New("Unexpected request type") }
        if v != nil {
            if err := v(hs); err != nil { return nil, err }
        }
        return &HandshakeRsp{}, nil
    }

    // Create the send-handshake-and-verify-acknowledged function
    // This function sends a HandshakeReq and expects a HandshakeRsp (ACK).
    asker := func() error {
        rsp, err := conn.Query(hs_local)
        if err != nil { return err }
        _, ok := rsp.(*HandshakeRsp)
        if ! ok { return nil }
        return errors.New("Unexpected response type")
    }

    // Run both routines
    ch := make(chan error)
    go func(){ ch <- conn.Handle(handler) }()
    go func(){ ch <- asker() }()

    // Get both responses
    err1 := <- ch
    err2 := <- ch

    // Decide which error is more important
    return conn, pick_error(err1, err2)
}

func pick_error(a, b error) error {
    if a == nil { return b }
    if b == nil { return a }
    if a == obp.ErrSocketClosed { return b }
    if b == obp.ErrSocketClosed { return a }
    return errors.New(fmt.Sprintf("Errors: %v, %v", a, b))
}

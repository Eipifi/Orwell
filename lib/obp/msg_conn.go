package obp
import "orwell/lib/butils"

type Msg butils.Chunk

type MsgDecoder func([]byte) (Msg, error)
type MsgEncoder func(Msg) ([]byte, error)
type MsgHandler func(Msg) (Msg, error)

type MsgConn struct {
    conn *Conn
    encoder MsgEncoder
    decoder MsgDecoder
}

func NewMsgConn(conn *Conn, encoder MsgEncoder, decoder MsgDecoder) *MsgConn {
    return &MsgConn{conn, encoder, decoder}
}

func (c *MsgConn) Close() {
    c.conn.Close()
}

func (c *MsgConn) Query(request Msg) (rsp Msg, err error) {
    buf, err := c.encoder(request)
    if err != nil { return }
    res, err := c.conn.Query(buf)
    if err != nil { return }
    return c.decoder(res)
}

func (c *MsgConn) Handle(handler MsgHandler) error {
    adapted_handler := func(req []byte) (rsp []byte, err error) {
        msg_req, err := c.decoder(req)
        if err != nil { return }
        msg_rsp, err := handler(msg_req)
        if err != nil { return }
        return c.encoder(msg_rsp)
    }
    return c.conn.Handle(adapted_handler)
}
package serv
import (
    "log"
    "orwell/lib/obp"
    "net"
    "orwell/lib/logging"
    "errors"
    "orwell/lib/protocol/orcache"
)

type Peer struct {
    conn *obp.MsgConn
    log *log.Logger
}

func Talk(socket net.Conn) {
    if err := TalkTo(socket); err != nil {
        log.Println(err)
    }
}

func TalkTo(socket net.Conn) (err error) {
    p := &Peer{}
    p.log = logging.GetLogger(socket.RemoteAddr().String())
    p.conn = orcache.Connection(socket)

    defer p.conn.Close()
    for {
        err = p.conn.Handle(p.messageHandler)
        if err != nil {
            break
        }
    }
    return nil
}

func (p *Peer) Close() {
    p.conn.Close()
}

func (p *Peer) messageHandler(msg obp.Msg) (rsp obp.Msg, err error) {
    return nil, errors.New("Unknown message type")
}

func (p *Peer) DoHandshake(local_hs *orcache.Handshake) (remote_hs *orcache.Handshake) {

}
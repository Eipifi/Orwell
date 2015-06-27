package main
import (
    "orwell/lib/obp"
    "net"
    "orwell/lib/protocol/orchain"
    "errors"
    "log"
    "orwell/lib/logging"
)

type Peer struct {
    conn *obp.MsgConn
    log *log.Logger
    hs *orchain.HandshakeReq
}

func Connect(socket net.Conn) {
    var err error
    p := &Peer{}
    p.log = logging.GetLogger(socket.RemoteAddr().String())
    if p.conn, err = orchain.Connect(socket, GenerateHandshake(), p.verifyHandshake); err != nil {
        p.conn.Close()
        p.log.Println(err)
        return
    }
    PeerJoined(p, p.hs)
    defer p.conn.Close()
    defer PeerLeft(p)
    for {
        if err = p.conn.Handle(p.handleRequest); err != nil {
            p.log.Println(err)
            return
        }
    }
}

func (p *Peer) verifyHandshake(hs *orchain.HandshakeReq) error {
    p.hs = hs
    return nil // we always accept the handshake
}

func (p *Peer) handleRequest(req obp.Msg) (rsp obp.Msg, err error) {
    p.log.Printf("Received request: %+v", req)
    return nil, errors.New("Not implemented")
}
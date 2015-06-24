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

func HandleConnection(socket net.Conn) (err error) {
    hs := &orchain.HandshakeReq{}
    p := &Peer{}
    p.log = logging.GetLogger(socket.RemoteAddr().String())
    if p.conn, err = orchain.Connect(socket, hs, nil); err != nil {
        p.log.Printf("Error: %v", err)
        return
    }
    defer p.conn.Close()
    for {
        if err = p.conn.Handle(p.handleRequest); err != nil {
            p.log.Printf("Error: %v", err)
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
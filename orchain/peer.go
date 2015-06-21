package main
import (
    "orwell/lib/obp"
    "net"
    "orwell/lib/protocol/orchain"
    "time"
)

type Peer struct {
    conn *obp.MsgConn
}

func NewPeer(socket net.Conn) *Peer {
    p := &Peer{orchain.Conn(socket)}
    return p
}

func (p *Peer) Lifecycle() {

    time.Sleep(time.Second * 5)
    p.conn.Close()
}

/*
    Handshake:
        Peer info exchange - supported version, capabilities, etc
        Two way, each side sends its own and awaits for response



*/
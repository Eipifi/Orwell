package main
import (
    "net"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/conv"
    "log"
    "os"
)

type Peer struct {
    cn net.Conn
    hs *orcache.Handshake
    log *log.Logger
}

func HandleConnection(conn net.Conn) {
    go func(){
        prefix := conn.RemoteAddr().String() + " "
        peer := &Peer{cn: conn, log: log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)}
        err := peer.lifecycle()
        if err != nil {
            peer.log.Println(err)
        }
    }()
}

func (p *Peer) lifecycle() (err error) {
    defer p.close()
    p.log.Println("Connected")
    if p.hs, err = conv.ShakeHands(p.cn, "orcache", nil); err != nil { return }
    p.log.Println("HS:", p.hs)

    inbox := conv.MessageListener(p.cn)
    for {
        select {
            case msg := <- inbox:
                if msg == nil { return }
                p.log.Println("Received", msg)
                p.handleMessage(msg)
        }
    }
    return
}

func (p *Peer) close() error {
    p.log.Println("Disconnected")
    return p.cn.Close()
}
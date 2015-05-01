package main
import (
    "net"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/conv"
    "log"
    "os"
    "orwell/orlib/butils"
)

type Peer struct {
    cn net.Conn
    hs *orcache.Handshake
    Log *log.Logger
    Out chan<- butils.Chunk
    GetOrders *RequestRouter
    PutOrders *RequestRouter
}

func HandleConnection(conn net.Conn) {
    go func(){
        prefix := conn.RemoteAddr().String() + " "
        peer := &Peer{cn: conn, Log: log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)}
        err := peer.lifecycle()
        if err != nil {
            peer.Log.Println(err)
        }
    }()
}

func (p *Peer) lifecycle() (err error) {
    defer p.close()
    if p.hs, err = conv.ShakeHands(p.cn, "orcache", nil); err != nil { return }
    inbox := conv.MessageListener(p.cn)
    p.Out = conv.MessageSender(p.cn)
    p.GetOrders = NewRouter(p.Out)
    p.PutOrders = NewRouter(p.Out)
    for {
        select {
            case msg := <- inbox:
                if msg == nil { return }
                p.handleMessage(msg)
        }
    }
    return
}

func (p *Peer) close() error {
    close(p.Out)
    p.GetOrders.Close()
    p.PutOrders.Close()
    return p.cn.Close()
}

func (p *Peer) AsyncSend(msg butils.Chunk) {
    go func(){ p.Out <- msg }()
}
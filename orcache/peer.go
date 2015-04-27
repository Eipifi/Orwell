package main
import (
    "net"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/conv"
    "log"
    "os"
    "orwell/orlib/butils"
    "orwell/orlib/protocol/common"
)

type Peer struct {
    cn net.Conn
    hs *orcache.Handshake
    Log *log.Logger
    Out chan<- butils.Chunk
    ToDo chan GetOrder
    Orders map[common.Token] GetOrder
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
    p.Log.Println("Connected")
    if p.hs, err = conv.ShakeHands(p.cn, "orcache", nil); err != nil { return }
    p.Log.Println("HS:", p.hs)

    inbox := conv.MessageListener(p.cn)
    p.Out = conv.MessageSender(p.cn)
    p.ToDo = make(chan GetOrder)
    for {
        select {
            case msg := <- inbox:
                if msg == nil { return }
                p.Log.Println("Received", msg)
                p.handleMessage(msg)
            case order := <- p.ToDo:
                p.handleOrder(order)
        }
    }
    return
}

func (p *Peer) close() error {
    p.Log.Println("Disconnected")
    close(p.Out)
    return p.cn.Close()
}
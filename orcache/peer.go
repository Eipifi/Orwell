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
    ToGet chan GetOrder
    ToPut chan PutOrder
    GetOrders map[common.Token] GetOrder
    PutOrders map[common.Token] PutOrder
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
    p.ToGet = make(chan GetOrder)
    p.ToPut = make(chan PutOrder)
    for {
        select {
            case msg := <- inbox:
                if msg == nil { return }
                p.Log.Println("Received", msg)
                p.handleMessage(msg)
            case order := <- p.ToGet: p.handleGetOrder(order)
            case order := <- p.ToPut: p.handlePutOrder(order)
        }
    }
    return
}

func (p *Peer) close() error {
    p.Log.Println("Disconnected")
    close(p.Out)
    return p.cn.Close()
}

func (p *Peer) AsyncSend(msg butils.Chunk) {
    p.Log.Println("sending", msg)
    go func(){ p.Out <- msg }()
}
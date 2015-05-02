package main
import (
    "net"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/conv"
    "log"
    "os"
    "orwell/orlib/butils"
    "sync"
)

type Peer struct {
    cn net.Conn
    Hs *orcache.Handshake
    Log *log.Logger
    out chan<- butils.Chunk
    GetOrders *RequestRouter
    PutOrders *RequestRouter
    closed bool
    mtx *sync.Mutex
}

func HandleConnection(conn net.Conn) {
    go func(){
        prefix := conn.RemoteAddr().String() + " "
        peer := &Peer{}
        peer.cn = conn
        peer.Log = log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
        peer.closed = false
        peer.mtx = &sync.Mutex{}
        err := peer.lifecycle()
        if err != nil {
            peer.Log.Println(err)
        }
    }()
}

func (p *Peer) lifecycle() (err error) {
    // Establish connection by exchanging handshakes
    if p.Hs, err = conv.ShakeHands(p.cn, "orcache", nil); err != nil { return }
    // Ensure proper finalization
    defer p.close()
    // Start listening on incoming messages (channel will close on socket close / error)
    inbox := conv.MessageListener(p.cn)
    // Start sending outgoing messages (socket will be closed on chan close / error)
    p.out = conv.MessageSender(p.cn)
    // Start the GET request manager
    p.GetOrders = NewRouter(p.out)
    // Start the PUT request manager
    p.PutOrders = NewRouter(p.out)
    // Finally announce peer presence to manager
    Manager.Join(p)
    // Loop
    for {
        select {
            case msg := <- inbox: // If new message was received
                if msg == nil { return } // If nil, inbox chan was closed. Peer must close.
                p.handleMessage(msg) // Else handle message
        }
    }
    return
}

// Socket closing procedure. Can be called multiple times.
func (p *Peer) close() error {
    p.mtx.Lock()
    defer p.mtx.Unlock()
    if p.closed { return nil }
    p.closed = true
    Manager.Leave(p)        // announce that the peer is no longer reachable
    close(p.out)            // close the send channel
    p.GetOrders.Close()     // cancel all get orders
    p.PutOrders.Close()     // cancel all put orders
    return p.cn.Close()     // finally close the socket
}

func (p *Peer) Send(msg butils.Chunk) {
    p.out <- msg
}

func (p *Peer) handleMessage(msg butils.Chunk) {
    switch msg := msg.(type) {
        case *orcache.GetReq:       go p.Send(Find(msg))
        case *orcache.GetRsp:       if !p.GetOrders.Respond(msg) { p.close() }
        case *orcache.PublishReq:   go p.Send(Publish(msg))
        case *orcache.PublishRsp:   if !p.PutOrders.Respond(msg) { p.close() }
        default: panic("Unrecognized message type")
    }
}
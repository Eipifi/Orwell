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
    Hs *orcache.Handshake
    GetOrders *RequestRouter
    PutOrders *RequestRouter

    cn net.Conn
    log *log.Logger
    out chan<- butils.Chunk
    closed bool
    mtx *sync.Mutex
    mgr Manager
}

func NewPeer(conn net.Conn, mgr Manager) *Peer {
    prefix := conn.RemoteAddr().String() + " "
    peer := &Peer{}
    peer.cn = conn
    peer.log = log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    peer.closed = false
    peer.mtx = &sync.Mutex{}
    peer.mgr = mgr
    go peer.lifecycle()
    return peer
}

func (p *Peer) lifecycle() {
    var err error
    // Establish connection by exchanging handshakes
    if p.Hs, err = conv.ShakeHands(p.cn, "orcache", p.mgr.LocalAddress()); err != nil {
        p.log.Println("Handshake exchange failed:", err)
        return
    }
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
    p.mgr.Join(p)
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
    p.mgr.Leave(p)        // announce that the peer is no longer reachable
    close(p.out)            // close the send channel
    p.GetOrders.Close()     // cancel all get orders
    p.PutOrders.Close()     // cancel all put orders
    return p.cn.Close()     // finally close the socket
}

// Sends the given message to remote peer
func (p *Peer) Send(msg butils.Chunk) {
    p.out <- msg
}

func (p *Peer) handleMessage(msg butils.Chunk) {
    switch msg := msg.(type) {
        case *orcache.GetReq:       go p.Send(Find(msg, p.mgr))
        case *orcache.GetRsp:       if !p.GetOrders.Respond(msg) { p.close() }
        case *orcache.PublishReq:   go p.Send(Publish(msg, p.mgr))
        case *orcache.PublishRsp:   if !p.PutOrders.Respond(msg) { p.close() }
        default: panic("Unrecognized message type")
    }
}
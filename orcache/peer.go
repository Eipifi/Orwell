package main
import (
    "net"
    "orwell/orlib/protocol/orcache"
    "log"
    "os"
    "sync"
    "orwell/orlib/client"
    "io"
)

type Peer struct {
    Hs *orcache.Handshake
    FetchOrders *RequestRouter
    PublishOrders *RequestRouter

    cn net.Conn
    log *log.Logger
    out chan<- orcache.Message
    closed bool
    mtx *sync.Mutex
    mgr *Manager
}

func NewPeer(conn net.Conn, mgr *Manager) *Peer {
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
    if p.Hs, err = client.ShakeHands(p.cn, "orcache", p.mgr.AdvertisedID, p.mgr.AdvertisedPort, p.mgr.CheckHandshake); err != nil {
        p.log.Println("Handshake exchange failed:", err)
        return
    }
    // Ensure proper finalization
    defer p.close()
    // Start listening on incoming messages (channel will close on socket close / error)
    inbox := MessageListener(p.cn)
    // Start sending outgoing messages (socket will be closed on chan close / error)
    p.out = MessageSender(p.cn)
    // Start the FETCH request manager
    p.FetchOrders = NewRouter(p.out)
    // Start the PUBLISH request manager
    p.PublishOrders = NewRouter(p.out)
    // Finally announce peer presence to manager
    p.mgr.Join(p)
    // Loop
    for {
        select {
            case msg := <- inbox: // If new message was received
                if msg == nil { return } // If nil, inbox chan was closed. Peer must close.
                go p.handleMessage(msg) // Else handle message
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
    p.mgr.Leave(p)          // announce that the peer is no longer reachable
    close(p.out)            // close the send channel
    p.FetchOrders.Close()   // cancel all get orders
    p.PublishOrders.Close() // cancel all put orders
    return p.cn.Close()     // finally close the socket
}

// Sends the given message to remote peer
func (p *Peer) Send(msg orcache.Message) {
    p.out <- msg
}

func (p *Peer) handleMessage(msg orcache.Message) {
    switch msg := msg.(type) {
        case *orcache.FetchReq:     p.Send(Find(msg, p.mgr))
        case *orcache.FetchRsp:     if !p.FetchOrders.Respond(msg) { p.close() }
        case *orcache.PublishReq:   p.Send(Publish(msg, p.mgr))
        case *orcache.PublishRsp:   if !p.PublishOrders.Respond(msg) { p.close() }
        case *orcache.PeersReq:     p.Send(&orcache.PeersRsp{p.mgr.FindAddresses(p.Hs.ID)})
        case *orcache.PeersRsp:     p.log.Println(msg)
        default: panic("Unrecognized message type")
    }
}

func MessageListener(conn io.Reader) <-chan orcache.Message {
    c := make(chan orcache.Message)
    go func(){
        defer close(c)
        for {
            msg, err := orcache.ReadAnyMessage(conn)
            if err != nil { return }
            c <- msg
        }
    }()
    return c
}

func MessageSender(conn io.WriteCloser) chan<- orcache.Message {
    c := make(chan orcache.Message)
    go func(){
        defer conn.Close()
        for {
            msg, ok := <- c
            if ! ok { return }
            if orcache.WriteMessage(conn, msg) != nil { return }
        }
    }()
    return c
}
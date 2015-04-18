package main
import (
    "net"
    "orwell/orlib/protocol"
)

const UserAgent = "Orcache"

type Peer struct {
    Upstream chan protocol.Msg
    NewJobs chan *GetJob

    jobs map[protocol.Token] *GetJob
    downstream <-chan protocol.Msg
    conn net.Conn
    hs *protocol.Handshake
    env *Env
    alive bool
}

func NewPeer(conn net.Conn, env *Env) (*Peer) {
    p := &Peer{}
    p.Upstream = make(chan protocol.Msg)
    p.NewJobs = make(chan *GetJob)
    p.conn = conn
    p.alive = true
    p.env = env
    return p
}

func (p *Peer) Lifecycle() {
    Info.Println("Connected to peer", p.conn.RemoteAddr())
    defer p.Close()
    if err := p.exchangeHandshakes(); err != nil { return }
    Info.Println("Successfully exchanged handshakes with", p.conn.RemoteAddr())
    p.downstream = readMessages(p.conn)
    for {
        if !p.alive { break }
        select {
            case job := <- p.NewJobs:
                p.handleJob(job)
            case msg := <- p.Upstream:
                p.sendMsg(msg)
            case msg := <- p.downstream:
                if msg == nil { return }
                p.handleMsg(msg)
        }
    }
}

func (p *Peer) Close() {
    if !p.alive { return }
    Info.Println("Closing connection with peer", p.conn.RemoteAddr())
    p.alive = false
    close(p.Upstream)
    close(p.NewJobs)
    p.conn.Close()
}

func (p *Peer) exchangeHandshakes() (err error) {
    // Initialize
    r := protocol.NewReader(p.conn)
    w := protocol.NewWriter()

    // Send our Handshake
    w.WriteFramedMessage(&protocol.Handshake{protocol.OrcacheMagic, protocol.SupportedVersion, UserAgent, nil})
    if err = w.Commit(p.conn); err != nil { return }

    // Await for the Handshake
    p.hs = &protocol.Handshake{}
    if err = r.ReadSpecificFramedMessage(p.hs); err != nil { return }

    // Send the HandshakeAck
    w.WriteFramedMessage(&protocol.HandshakeAck{})
    if err = w.Commit(p.conn); err != nil { return }

    // Await for the HandshakeAck
    var ack protocol.HandshakeAck
    if err = r.ReadSpecificFramedMessage(&ack); err != nil { return }
    return
}

func (p *Peer) handleJob(job *GetJob) {
    Info.Println("Received job", job, "for peer", p.conn.RemoteAddr())
    // So we received a job. Let's look at it.
    if _, ok := p.jobs[job.Msg.Token]; ok {
        // Hmm. It looks like we are already dealing with a job of this ID.
        // We cannot interfere with it. Therefore we quit the job immediately.
        job.Fail()
    } else {
        // Ok, this token is unknown. We can proceed.
        p.jobs[job.Msg.Token] = job
        p.sendMsg(job.Msg)
    }
}

func (p *Peer) sendMsg(msg protocol.Msg) {
    // todo: assess the connection state
    Info.Println("Sent message", msg, "to peer", p.conn.RemoteAddr())
    w := protocol.NewWriter()
    w.WriteFramedMessage(msg)
    w.Commit(p.conn)
}

func (p *Peer) handleMsg(msg protocol.Msg) {
    Info.Println("Received message", msg, "from peer", p.conn.RemoteAddr())
    // Switch over message type
    switch msg := msg.(type) {
        case *protocol.Get:
        p.HandleGet(msg)

        case *protocol.CardFound:
        p.HandleCardFound(msg)

        case *protocol.CardNotFound:
        p.HandleCardNotFound(msg)

        default:
        panic("Unrecognized message type")
    }
}

func (p *Peer) MaybeSendMsg(msg protocol.Msg) bool {
    return Maybe(func(){
        p.Upstream <- msg
    })
}

func (p *Peer) MaybeSendJob(job *GetJob) bool {
    return Maybe(func(){
        p.NewJobs <- job
    })
}

func (p *Peer) HandleGet(msg *protocol.Get) {
    go func(){
        result := Find(msg, p.env)
        if result.Bytes == nil {
            Info.Println("Card", msg.ID, "not found, responsing to peer", p.conn.RemoteAddr())
            p.MaybeSendMsg(&protocol.CardNotFound{msg.Token, result.TTL})
        } else {
            Info.Println("Card", msg.ID, "found, responsing to peer", p.conn.RemoteAddr())
            p.MaybeSendMsg(&protocol.CardFound{msg.Token, result.Bytes})
        }
    }()
}

func (p *Peer) HandleCardFound(msg *protocol.CardFound) {
    p.completeJob(msg.Token, &GetResponse{msg.Card, 0})
}

func (p *Peer) HandleCardNotFound(msg *protocol.CardNotFound) {
    p.completeJob(msg.Token, &GetResponse{nil, msg.TTL})
}

func (p *Peer) completeJob(token protocol.Token, resp *GetResponse) {
    // Fetch the specified job
    job, ok := p.jobs[token]
    // Do we have it?
    if !ok {
        p.Close()
        return
    }
    // Send the response and clean up
    job.Sink <- resp
    delete(p.jobs, token)
}

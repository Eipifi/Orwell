package main
import (
    "net"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/comm"
    "orwell/orlib/protocol/types"
    "fmt"
)

const UserAgent = "Orcache"

type Peer struct {
    ms *orcache.OrcacheMessenger
    Upstream chan comm.Msg
    NewJobs chan *GetJob
    jobs map[types.Token] *GetJob
    downstream chan comm.Msg
    env *Env
    alive bool
}

func NewPeer(env *Env) (*Peer) {
    p := &Peer{}
    p.Upstream = make(chan comm.Msg)
    p.NewJobs = make(chan *GetJob)
    p.downstream = make(chan comm.Msg)
    p.alive = true
    p.env = env
    return p
}

func (p *Peer) Lifecycle(conn net.Conn) {
    var err error
    p.ms, err = orcache.NewOrcacheMessenger(conn, UserAgent, nil)
    if err != nil {

    }
    defer p.Close()

    go func(){
        for {
            m, e := p.ms.ReadAny()
            if e != nil {
                fmt.Println(e)
                break
            }
            p.downstream <- m
        }
        close(p.downstream)
    }()

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
    p.alive = false
    close(p.Upstream)
    close(p.NewJobs)
    p.ms.Close()
}

func (p *Peer) handleJob(job *GetJob) {
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

func (p *Peer) sendMsg(msg comm.Msg) {
    // todo: assess the connection state
    p.ms.Write(msg)
}

func (p *Peer) handleMsg(msg comm.Msg) {
    // Switch over message type
    switch msg := msg.(type) {
        case *orcache.Get:
        p.HandleGet(msg)

        case *orcache.CardFound:
        p.HandleCardFound(msg)

        case *orcache.CardNotFound:
        p.HandleCardNotFound(msg)

        default:
        panic("Unrecognized message type")
    }
}

func (p *Peer) MaybeSendMsg(msg comm.Msg) bool {
    return Maybe(func(){
        p.Upstream <- msg
    })
}

func (p *Peer) MaybeSendJob(job *GetJob) bool {
    return Maybe(func(){
        p.NewJobs <- job
    })
}

func (p *Peer) HandleGet(msg *orcache.Get) {
    go func(){
        result := Find(msg, p.env)
        if result.Bytes == nil {
            p.MaybeSendMsg(&orcache.CardNotFound{msg.Token, result.TTL})
        } else {
            p.MaybeSendMsg(&orcache.CardFound{msg.Token, result.Bytes})
        }
    }()
}

func (p *Peer) HandleCardFound(msg *orcache.CardFound) {
    p.completeJob(msg.Token, &GetResponse{msg.Card, 0})
}

func (p *Peer) HandleCardNotFound(msg *orcache.CardNotFound) {
    p.completeJob(msg.Token, &GetResponse{nil, msg.TTL})
}

func (p *Peer) completeJob(token types.Token, resp *GetResponse) {
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
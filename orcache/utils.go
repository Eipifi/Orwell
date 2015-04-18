package main
import (
    "orwell/orlib/protocol"
    "io"
)

func Maybe(routine func()) bool {
    ret := true
    defer func(){
        ret = (recover() == nil)
    }()
    routine()
    return ret
}

func readMessages(source io.Reader) <-chan protocol.Msg{
    c := make(chan protocol.Msg) // maybe introduce buffered chan?
    go func(){
        r := protocol.NewReader(source)
        for {
            msg, err := r.ReadFramedMessage()
            if err != nil { break }
            c <- msg
        }
        close(c)
    }()
    return c
}

type GetJob struct {
    Msg *protocol.Get
    Sink chan *GetResponse
}

type GetResponse struct {
    Bytes []byte
    TTL protocol.TTL
}


func (j *GetJob) Cancel() {
    // Job did not get to a peer
    j.Sink <- nil
}

func (j *GetJob) Fail() {
    // Response was not received
    j.Sink <- &GetResponse{nil, j.Msg.TTL}
}

type Env struct {
    Manager *Manager
    Cache Cache
    Locker TokenLocker
}
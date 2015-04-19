package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/types"
)

func Maybe(routine func()) bool {
    ret := true
    defer func(){
        ret = (recover() == nil)
    }()
    routine()
    return ret
}

type GetJob struct {
    Msg *orcache.Get
    Sink chan *GetResponse
}

type GetResponse struct {
    Bytes []byte
    TTL types.TTL
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
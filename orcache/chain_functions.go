package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/crypto/hash"
)

func GetRspValidator(req *orcache.FetchReq) func(orcache.TokenMessage) bool {
    return func(msg orcache.TokenMessage) bool {
        switch msg := msg.(type) {
            case *orcache.FetchRsp:
                if msg.Card == nil { return true }
                return hash.Equal(msg.Card.Key.Id(), req.ID) && uint64(msg.Card.Payload.Version) >= req.Version // TODO: fix this damn int64
            default: return false
        }
    }
}

func Find(req *orcache.FetchReq, pf PeerFinder) (rsp *orcache.FetchRsp) {
    rsp = &orcache.FetchRsp{req.Token, req.TTL, nil}
    if rsp.Card = Storage.Get(req.ID, req.Version); rsp.Card != nil { return }
    if Locker.Lock(rsp.Token) {
        validator := GetRspValidator(req)
        defer Locker.Unlock(rsp.Token)
        for {
            if rsp.TTL == 0 { return }
            peer := pf.FindPeer(req.ID)
            if peer == nil { return }
            rsp.TTL -= 1
            if r := peer.FetchOrders.Ask(&orcache.FetchReq{rsp.Token, rsp.TTL, req.ID, req.Version}, validator); r != nil {
                nrsp := r.(*orcache.FetchRsp) // correct type assumed
                if nrsp.TTL < rsp.TTL { rsp.TTL = nrsp.TTL }
                if nrsp.Card != nil {
                    rsp.Card = nrsp.Card
                    return
                }
            }
            if rsp.Card = Storage.Get(req.ID, req.Version); rsp.Card != nil { return }
        }
    }
    return
}

func Publish(req *orcache.PublishReq, pf PeerFinder) (rsp *orcache.PublishRsp) {
    rsp = &orcache.PublishRsp{req.Token, req.TTL}
    Storage.Put(req.Card)
    if Locker.Lock(req.Token) {
        defer Locker.Unlock(req.Token)
        for {
            if rsp.TTL == 0 { return }
            peer := pf.FindPeer(req.Card.Key.Id())
            if peer == nil { return }
            if r := peer.PublishOrders.Ask(&orcache.PublishReq{rsp.Token, rsp.TTL, req.Card}, nil); r != nil {
                nrsp := r.(*orcache.PublishRsp) // correct type assumed
                if nrsp.TTL < rsp.TTL { rsp.TTL = nrsp.TTL }
            }
        }
    }
    return
}
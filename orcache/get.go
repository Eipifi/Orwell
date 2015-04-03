package main
import (
    "orwell/orlib/protocol"
)

type GetResponse struct {
    Card []byte
    TTL uint8
}

func getHandler(req *protocol.Get, cancel <-chan bool, env Env) *GetResponse {
    if data := env.Cache.Get(req.ID); data != nil {
        return &GetResponse{data, 0}
    }
    ttl := req.TTL
    if ttl > 0 && env.Locker.Lock(req.Token) {
        defer env.Locker.Unlock(req.Token)
        for ttl > 0 {
            ttl -= 1
            peer := env.Manager.PickPeer(req.ID)
            if peer == nil { break }
            select {
                case <- cancel:
                    return nil
                case rsp := <- peer.RelayGet(&protocol.Get{req.Token, ttl, req.ID, req.Version}):
                    if rsp != nil {
                        if rsp.Card != nil {
                            // TODO: validate card (check if matches the given ID/version)
                            // TODO: put the card in cache
                            return &GetResponse{rsp.Card}, 0
                        } else if rsp.TTL < ttl {
                            ttl = rsp.TTL
                        }
                    }
            }
            if data := env.Cache.Get(req.ID); data != nil {
                return &GetResponse{data, 0}
            }
        }
    }
    return &GetResponse{nil, 0}
}
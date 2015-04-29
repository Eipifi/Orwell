package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/common"
    "orwell/orlib/crypto/card"
)

type GetResult struct {
    Card *card.Card
    TTL common.TTL
}

func Find(req *orcache.Get) (res GetResult) {
    res.TTL = req.TTL
    if res.Card = Storage.Get(req.ID, req.Version); res.Card != nil { return }
    if Locker.Lock(req.Token) {
        defer Locker.Unlock(req.Token)
        for {
            peer := FindPeer(req.ID)
            if peer == nil { return }
            res.TTL -= 1
            if r := peer.AskGet(&orcache.Get{req.Token, res.TTL, req.ID, req.Version}); r != nil {
                if r.Card != nil {
                    Storage.Put(r.Card)
                    return *r
                }
                if r.TTL < res.TTL { res.TTL = r.TTL }
            }
            if res.Card = Storage.Get(req.ID, req.Version); res.Card != nil { return }
        }
    }
    return
}
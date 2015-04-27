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
    if res.Card = Storage.Get(req.ID, req.Version); res.Card != nil { return }
    if Locker.Lock(req.Token) {
        defer Locker.Unlock(req.Token)
        for res.TTL = req.TTL - 1; res.TTL > 0; res.TTL-- {
            peer := FindPeer(req.ID)
            if peer == nil { return }
            if r := peer.AskGet(&orcache.Get{req.Token, res.TTL, req.ID, req.Version}); r != nil {
                if r.Card != nil {
                    Storage.Put(req.ID, r.Card)
                    return *r
                }
                if r.TTL < res.TTL { res.TTL = r.TTL }
            }
            if res.Card = Storage.Get(req.ID, req.Version); res.Card != nil { return }
        }
    } else { res.TTL = req.TTL }
    return
}
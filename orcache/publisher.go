package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/common"
)

func Publish(msg *orcache.Publish) (ttl common.TTL) {
    ttl = msg.TTL
    Storage.Put(msg.Card)
    if Locker.Lock(msg.Token) {
        defer Locker.Unlock(msg.Token)
        for {
            peer := FindPeer(msg.Card.Key.Id())
            if peer == nil { return }
            ttl -= 1
            result := peer.AskPut(&orcache.Publish{msg.Token, ttl, msg.Card})
            if result != nil {
                if *result < ttl { ttl = *result }
            }
        }
    }
    return
}
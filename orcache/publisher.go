package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/types"
    "orwell/orlib/crypto/card"
)

func Publish(msg *orcache.Publish, env *Env) (ttl types.TTL, err error) {
    c := &card.Card{}
    if err = c.UnmarshalBinary(msg.Card); err != nil { return }
    // todo: implement chain publishing
    env.Cache.Put(c.PubKey().Id(), msg.Card)
    ttl = msg.TTL
    return
}
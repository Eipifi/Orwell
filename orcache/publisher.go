package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/types"
    "orwell/orlib/card"
)

func Publish(msg *orcache.Publish, env *Env) (ttl types.TTL, err error) {
    var c *card.Card
    if c, err = card.Unmarshal(msg.Card); err != nil { return }
    // todo: implement chain publishing
    env.Cache.Put(c.Key.Id(), msg.Card)
    ttl = msg.TTL
    return
}
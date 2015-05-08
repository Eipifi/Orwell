package main
import (
    "orwell/orlib/crypto/hash"
    "orwell/orlib/crypto/card"
)

/*
    This is obviously a naive implementation.
    1. Implement LRU
    2. Provide concurrency safety
*/

var Storage Cache = NewMapCache()

type Cache interface {
    Get(hash.ID, uint64) *card.Card
    Put(*card.Card)
}

type MapCache struct {
    data map[hash.ID] *card.Card
}

func NewMapCache() *MapCache {
    m := &MapCache{}
    m.data = make(map[hash.ID] *card.Card)
    return m
}

func (c *MapCache) Get(id hash.ID, version uint64) *card.Card {
    return c.data[id]
}

func (c *MapCache) Put(val *card.Card) {
    c.data[val.Key.Id()] = val
}
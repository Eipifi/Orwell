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

var Storage Cache = &MapCache{}

type Cache interface {
    Get(*hash.ID, uint64) *card.Card
    Put(*hash.ID, *card.Card)
}

type MapCache struct {
    data map[hash.ID] *card.Card
}

func (c *MapCache) Get(id *hash.ID, version uint64) *card.Card {
    return c.data[*id]
}

func (c *MapCache) Put(id *hash.ID, val *card.Card) {
    c.data[*id] = val
}
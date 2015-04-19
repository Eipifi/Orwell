package main
import (
    "orwell/orlib/protocol/types"
    "fmt"
)

type Cache interface {
    Get(key *types.ID) []byte
    Put(key *types.ID, value []byte)
}

type EmptyCache struct { }

func (c *EmptyCache) Get(*types.ID) []byte {
    fmt.Println("Fetched!")
    return nil
}

func (c *EmptyCache) Put(*types.ID, []byte) {
    fmt.Println("Stored!")
}

type MapCache struct {
    data map[types.ID] []byte
}

func NewMapCache() *MapCache {
    c := &MapCache{}
    c.data = make(map[types.ID] []byte)
    return c
}

func (c *MapCache) Get(id *types.ID) []byte {
    return c.data[*id]
}

func (c *MapCache) Put(id *types.ID, val []byte) {
    c.data[*id] = val
}
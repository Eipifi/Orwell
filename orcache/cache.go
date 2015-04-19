package main
import (
    "orwell/orlib/protocol/types"
)

type Cache interface {
    Get(key *types.ID) []byte
    Put(key *types.ID, value []byte)
}

type EmptyCache struct { }

func (c *EmptyCache) Get(*types.ID) []byte { return nil }

func (c *EmptyCache) Put(*types.ID, []byte) { }
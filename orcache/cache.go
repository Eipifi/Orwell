package main
import "orwell/orlib/sig"

type Cache interface {
    Get(key sig.ID) []byte
    Put(key sig.ID, value []byte)
}

type EmptyCache struct { }

func (c *EmptyCache) Get(sig.ID) []byte { return nil }

func (c *EmptyCache) Put(sig.ID, []byte) { }
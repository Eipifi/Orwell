package main
import "orwell/orlib/sig"

type Cache interface {
    Get(key sig.ID) []byte
    Put(key sig.ID, value []byte)
}
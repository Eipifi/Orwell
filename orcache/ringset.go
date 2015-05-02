package main
import "orwell/orlib/crypto/hash"

type RingSet interface {
    Put(key *hash.ID, value interface{})
    GetClosest(key *hash.ID, max int) []interface{}
    Remove(key *hash.ID)
}



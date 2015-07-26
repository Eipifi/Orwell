package db
import (
    "github.com/boltdb/bolt"
)

var BUCKET_INFO = []byte("info")
var KEY_STATE = []byte("state")

func GetState(t *bolt.Tx) (s *State) {
    s = &State{}
    if Read(t, BUCKET_INFO, KEY_STATE, s) { return }
    return nil
}

func PutState(t *bolt.Tx, s *State) {
    Write(t, BUCKET_INFO, KEY_STATE, s)
}
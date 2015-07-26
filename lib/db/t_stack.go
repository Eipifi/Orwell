package db
import (
    "github.com/boltdb/bolt"
    "orwell/lib/protocol/orchain"
)

func PushBlock(t *bolt.Tx, block *orchain.Block) error {
    if err := VerifyNextBlock(t, block); err != nil { return err }
    s := GetState(t)
    PutBlock(t, block, s.Length)
    s.Length += 1
    s.Head = block.Header.ID()
    s.Work = s.Work.Plus(block.Header.Difficulty)
    PutState(t, s)
    return nil
}

func PopBlock(t *bolt.Tx) {
    s := GetState(t)
    if s.Length == 1 { panic("Can not remove the genesis block") }
    h := GetHeaderByID(t, s.Head)
    DelBlock(t, s.Head)
    s.Length -= 1
    s.Head = h.Previous
    s.Work = s.Work.Minus(h.Difficulty)
    PutState(t, s)
}
package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "github.com/boltdb/bolt"
)

var BUCKET_TXN = []byte("txn")

func PutTransaction(t *bolt.Tx, txn *orchain.Transaction) {
    tid := txn.ID()
    Write(t, BUCKET_TXN, tid[:], txn)
    for _, inp := range txn.Inputs {
        SetBillStatus(t, &inp, SPENT)
    }
    for i, _ := range txn.Outputs {
        SetBillStatus(t, &orchain.BillNumber{tid, uint64(i)}, UNSPENT)
    }
}

func GetTransaction(t *bolt.Tx, id foo.U256) *orchain.Transaction {
    res := &orchain.Transaction{}
    if Read(t, BUCKET_TXN, id[:], res) { return res }
    return nil
}

func DelTransaction(t *bolt.Tx, id foo.U256) {
    txn := GetTransaction(t, id)
    if txn == nil { return }
    tid := txn.ID()
    Del(t, BUCKET_TXN, id[:])
    for _, inp := range txn.Inputs {
        SetBillStatus(t, &inp, UNSPENT)
    }
    for i, _ := range txn.Outputs {
        SetBillStatus(t, &orchain.BillNumber{tid, uint64(i)}, NONEXISTENT)
    }
}
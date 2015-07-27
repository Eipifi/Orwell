package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
)

var BUCKET_TXN = []byte("txn")

func (t *Tx) PutTransaction(txn *orchain.Transaction) {
    tid := txn.ID()
    t.Write(BUCKET_TXN, tid[:], txn)
    for _, inp := range txn.Inputs {
        t.SetBillStatus(&inp, SPENT)
    }
    for i, _ := range txn.Outputs {
        t.SetBillStatus(&orchain.BillNumber{tid, uint64(i)}, UNSPENT)
    }
}

func (t *Tx) GetTransaction(id foo.U256) *orchain.Transaction {
    res := &orchain.Transaction{}
    if t.Read(BUCKET_TXN, id[:], res) { return res }
    return nil
}

func (t *Tx) DelTransaction(id foo.U256) {
    txn := t.GetTransaction(id)
    if txn == nil { return }
    tid := txn.ID()
    t.Del(BUCKET_TXN, id[:])
    for _, inp := range txn.Inputs {
        t.SetBillStatus(&inp, UNSPENT)
    }
    for i, _ := range txn.Outputs {
        t.SetBillStatus(&orchain.BillNumber{tid, uint64(i)}, NONEXISTENT)
    }
}
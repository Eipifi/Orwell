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
    if txn.Payload.Ticket != nil {
        t.PutTicket(*txn.Payload.Ticket)
    }
    if txn.Payload.Domain != nil {
        t.PutDomain(txn.Payload.Domain)
    }
}

func (t *Tx) GetTransaction(id foo.U256) *orchain.Transaction {
    res := &orchain.Transaction{}
    if t.Read(BUCKET_TXN, id[:], res) { return res }
    return nil
}
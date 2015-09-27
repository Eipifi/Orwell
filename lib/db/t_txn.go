package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "bytes"
    "io"
    "orwell/lib/utils"
)

var BUCKET_TXN = []byte("txn")
var BUCKET_TXN_LIST   = []byte("txn_list")

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

func (t *Tx) GetTransactionsFromBlock(id foo.U256) (txns []orchain.Transaction) {
    tids := t.Get(BUCKET_TXN_LIST, id[:])
    buf := bytes.NewBuffer(tids)
    for {
        var tid foo.U256
        err := tid.Read(buf)
        if err == io.EOF { break }
        utils.Ensure(err)
        txn := t.GetTransaction(tid)
        txns = append(txns, *txn)
    }
    return
}

func (t *Tx) PutTransactionsFromBlock(id foo.U256, txns []orchain.Transaction) {
    buf := &bytes.Buffer{}
    for _, txn := range txns {
        tid := txn.ID()
        tid.Write(buf)
        t.PutTransaction(&txn)
    }
    t.Put(BUCKET_TXN_LIST, id[:], buf.Bytes())
}
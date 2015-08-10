package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/utils"
)

func (t *Tx) ComputeBalance(txn *orchain.Transaction) (input, output uint64) {
    // TODO: check for overflows
    for _, inp := range txn.Inputs {
        bill := t.GetBill(&inp)
        input += bill.Value
    }
    for _, out := range txn.Outputs {
        output += out.Value
    }
    return
}

func (t *Tx) ComputeTransactionFee(txn *orchain.Transaction) (uint64, error) {
    in, out := t.ComputeBalance(txn)
    utils.Assert(in >= out)
    return in - out, nil
}
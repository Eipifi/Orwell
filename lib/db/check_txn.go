package db
import (
    "orwell/lib/protocol/orchain"
    "errors"
    "github.com/deckarep/golang-set"
)

func CheckTxnIsNew(t *Tx, txn *orchain.Transaction, is_first bool) error {
    // Check if no other transaction has the same ID
    // Note: https://github.com/bitcoin/bips/blob/master/bip-0030.mediawiki
    if t.GetTransaction(txn.ID()) != nil { return errors.New("Transaction ID already in use") }
    return nil
}

func CheckTxnProof(t *Tx, txn *orchain.Transaction, is_first bool) error {
    if is_first {
        if txn.Proof != nil { return errors.New("The proof is not allowed in a coinbase transaction") }
    } else {
        if txn.Proof == nil { return errors.New("The transaction lacks a proof") }
        if err := txn.Verify(); err != nil { return err }
        sender_id, err := txn.Proof.PublicKey.ID()
        if err != nil { return err }
        for _, inp := range txn.Inputs {
            bill := t.GetBill(&inp)
            if bill.Target != sender_id { return errors.New("Transaction tries to spend somebody else's bill") }
        }
    }
    return nil
}

func CheckTxnInputsUnspent(t *Tx, txn *orchain.Transaction, is_first bool) error {
    if is_first {
        if len(txn.Inputs) != 0 { return errors.New("Coinbase transaction can have no inputs") }
    }
    for _, inp := range txn.Inputs {
        if t.GetBillStatus(&inp) != UNSPENT { return errors.New("Input bill is already spent or does not exist") }
    }
    return nil
}

func CheckTxnNoDoubleSpend(t *Tx, txn *orchain.Transaction, is_first bool) error {
    to_spend := mapset.NewSet()
    for _, inp := range txn.Inputs {
        if ! to_spend.Add(inp) { return errors.New("Bill spent twice") }
    }
    return nil
}

func CheckTxnBalance(t *Tx, txn *orchain.Transaction, is_first bool) (err error) {
    if ! is_first {
        input, output := t.ComputeTxnInpOut(txn)
        if output > input { return errors.New("Transaction output must not be greater then its input") }
    }
    return
}
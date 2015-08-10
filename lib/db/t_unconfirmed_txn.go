package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
)

var BUCKET_TXN_UNCONFIRMED = []byte("txn_unconfirmed")

func (t *Tx) PutUnconfirmedTransaction(txn *orchain.Transaction) {
    /*
    tid := txn.ID()
    t.Write(BUCKET_TXN_UNCONFIRMED, tid[:], txn)
    */
}

func (t *Tx) DelUnconfirmedTransaction(tid foo.U256) {
    /*
    t.Del(BUCKET_TXN_UNCONFIRMED, tid[:])
    */
}

func (t *Tx) UnconfirmedTransactions() (result []orchain.Transaction) {
    /*
    b := t.tx.Bucket(BUCKET_TXN_UNCONFIRMED)
    var txn orchain.Transaction
    utils.Ensure(b.ForEach(func(k, v []byte) (err error) {
        err = butils.ReadAllInto(&txn, v)
        if err != nil { return }
        result = append(result, txn)
        return
    }))
    */
    return
}

func (t *Tx) RefreshUnconfirmedTransactions() {
    /*
    txns := t.UnconfirmedTransactions()
    stored := make([]orchain.Transaction, 0)
    t.DeleteAll(BUCKET_TXN_UNCONFIRMED)
    for _, txn := range txns {
        _, _, err := t.VerifyTransaction(&txn, false)
        if err == nil && t.VerifyTransactionDoesNotConflict(&txn, stored) == nil {
            t.PutUnconfirmedTransaction(&txn)
            stored = append(stored, txn)
        }
    }
    */
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

/*
    - store unconfirmed transactions,
    - make sure only non-conflicting txns are stored,
    - on conflict, remove less profitable ones (TODO)
*/

func (t *Tx) MaybeStoreUnconfirmedTransaction(txn *orchain.Transaction) (err error) {
    /*
    _, _, err = t.VerifyTransaction(txn, false)
    if err != nil { return }
    if err = t.VerifyTransactionDoesNotConflict(txn, t.UnconfirmedTransactions()); err != nil { return }
    t.PutUnconfirmedTransaction(txn)
    */
    return nil
}
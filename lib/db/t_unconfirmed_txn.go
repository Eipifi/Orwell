package db
import (
    "orwell/lib/protocol/orchain"
)

var BUCKET_TXN_UNCONFIRMED = []byte("txn_unconfirmed")

func (t *Tx) UnconfirmedTransactions() (result []orchain.Transaction) {
    return
}

func (t *Tx) RefreshUnconfirmedTransactions() {
}


func (t *Tx) MaybeStoreUnconfirmedTransaction(txn *orchain.Transaction) (err error) {
    return nil
}
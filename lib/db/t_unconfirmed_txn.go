package db
import (
    "orwell/lib/protocol/orchain"
    "github.com/deckarep/golang-set"
    "orwell/lib/utils"
    "errors"
    "orwell/lib/butils"
    "orwell/lib/foo"
)

var BUCKET_TXN_UNCONFIRMED = []byte("txn_unconfirmed")

func (t *Tx) UnconfirmedTransactions() (result []orchain.Transaction) {
    t.tx.Bucket(BUCKET_TXN_UNCONFIRMED).ForEach(func(k, v []byte) error {
        txn := orchain.Transaction{}
        utils.Ensure(butils.ReadAllInto(&txn, v))
        result = append(result, txn)
        return nil
    })
    return
}

func (t *Tx) RefreshUnconfirmedTransactions() {
    // we're lazy, we just remove all and reinsert. Because why not, we'll optimize later.
    current := t.UnconfirmedTransactions()
    c := t.tx.Bucket(BUCKET_TXN_UNCONFIRMED).Cursor()
    for k, _ := c.First(); k != nil; k, _ = c.Next() {
        utils.Ensure(c.Delete())
    }
    for _, txn := range current {
        t.MaybeStoreUnconfirmedTransaction(&txn) // we ignore failed inserts, that's the whole point of this operation.
    }
}


func (t *Tx) MaybeStoreUnconfirmedTransaction(txn *orchain.Transaction) (err error) {

    for _, check := range txn_checks {
        if err = check(t, txn, false); err != nil { return }
    }

    current := t.UnconfirmedTransactions()

    // Construct a set of inputs to be spent
    to_spend := mapset.NewSet()
    for _, cur := range current {
        for _, inp := range cur.Inputs {
            utils.Assert(to_spend.Add(inp))
        }
    }

    // Check if the new transaction tries to spend the same bill
    for _, inp := range txn.Inputs {
        if to_spend.Contains(inp) {
            return errors.New("The proposed transaction tries to spend a bill that other candidate transaction spends")
        }
    }

    // If the transaction carries a ticket
    if txn.Payload.Ticket != nil {
        // Count the tickets and find the cheapest one
        var tickets uint64 = 0
        var cheapest_id int = -1
        var cheapest_fee uint64 = foo.U64_MAX // max uint64
        for i, cur := range current {
            if cur.Payload.Ticket != nil {
                tickets += 1
                fee := t.ComputeFee(&cur)
                if fee < cheapest_fee {
                    cheapest_fee = fee
                    cheapest_id = i
                }
            }
        }

        allowed := t.AllowedTickets()

        if tickets > allowed {
            panic("This should never happen - we can not have more tickets in the unconfirmed buffer then allowed")
        }

        if tickets == allowed {
            // remove the cheapest txn
            tid := current[cheapest_id].ID()
            t.RawDel(BUCKET_TXN_UNCONFIRMED, tid[:])
        }
    }
    // add the new txn
    tid := txn.ID()
    t.RawWrite(BUCKET_TXN_UNCONFIRMED, tid[:], txn)
    return nil
}
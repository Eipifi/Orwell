package db
import (
    "orwell/lib/protocol/orchain"
    "github.com/deckarep/golang-set"
    "errors"
)

func CheckTxnsNoDuplicateTxns(t *Tx, txns []orchain.Transaction) error {
    txn_ids := mapset.NewSet()
    for _, txn := range txns {
        if ! txn_ids.Add(txn.ID()) { return errors.New("Duplicate transactions in block") }
    }
    return nil
}

func CheckTxnsNoDoubleSpend(t *Tx, txns []orchain.Transaction) error {
    // We'll collect all inputs and check for duplicates
    to_spend := mapset.NewSet()
    // For each transaction
    for _, txn := range txns {
        for _, inp := range txn.Inputs {
            if ! to_spend.Add(inp) { return errors.New("Two transactions in a block spend the same bill") }
        }
    }
    return nil
}

func CheckTxnsBalance(t *Tx, txns []orchain.Transaction) (err error) {
    // Here we'll store the sum of all fees
    var total_input_sum, total_output_sum uint64

    for _, txn := range txns {
        // TODO: check for overflows
        inp, out := t.ComputeTxnInpOut(&txn)
        total_input_sum += inp
        total_output_sum += out
    }

    // The transaction should also generate a reward
    state := t.GetState()
    total_input_sum += orchain.GetReward(state.Length)

    // Check if the sums match up
    if total_input_sum != total_output_sum {
        return errors.New("Invalid reward/fees")
    }

    return nil
}
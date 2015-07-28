package db
import (
    "orwell/lib/protocol/orchain"
    "errors"
    "orwell/lib/foo"
    "github.com/deckarep/golang-set"
    "orwell/lib/utils"
)

func (t *Tx) VerifyNextBlock(b *orchain.Block) (err error) {
    bid := b.Header.ID()
    state := t.GetState()

    // TODO: check timestamp, define the timestamp acceptance policy

    // Check if no other block has the same id
    // Note: while the chance of this happening is astronomically low, we still check this.
    // SHA256 might get broken at some point in the future.
    if t.GetHeaderByID(bid) != nil {
        return errors.New("Block with this ID already exists")
    }

    // Check if block is a child of the current head
    if ! foo.Equal(state.Head, b.Header.Previous) {
        return errors.New("The 'Previous' field of the block does not match the stored head")
    }

    // Check if the difficulty value is correct
    if t.GetDifficulty() != b.Header.Difficulty {
        return errors.New("Invalid difficulty value")
    }

    // Check if the block hash meets the specified difficulty
    if ! orchain.HashMeetsDifficulty(bid, b.Header.Difficulty) {
        return errors.New("Block hash does not meet the specified difficulty")
    }

    // Check if the timestamp is correct
    if state.Length > 0 {
        previous_header := t.GetHeaderByID(state.Head)
        if b.Header.Timestamp < previous_header.Timestamp {
            return errors.New("Block timestamp is smaller than the previous one")
        }
    }

    // Check if the Merkle root matches (and also if there is at least one transaction)
    if err = b.CheckMerkleRoot(); err != nil { return }

    // We'll collect all inputs and check for duplicates
    to_spend := mapset.NewSet()

    // We also check if there are no transaction duplicates
    txn_ids := mapset.NewSet()

    // Here we'll store the sum of all fees
    var total_input_sum, total_output_sum uint64

    // For each transaction
    for txn_num, txn := range b.Transactions {

        var tid foo.U256
        if tid, err = txn.TryID(); err != nil { return }

        // Check if no other transaction has the same ID
        // Note: https://github.com/bitcoin/bips/blob/master/bip-0030.mediawiki
        if t.GetTransaction(tid) != nil { return errors.New("Transaction ID already in use") }
        if ! txn_ids.Add(tid) { return errors.New("Duplicate transactions in block") }

        if txn_num == 0 { // Check the coinbase transaction
            if txn.Proof != nil { return errors.New("The proof is not required/allowed in a coinbase transaction") }
            if len(txn.Inputs) != 0 { return errors.New("Coinbase transaction can have no inputs") }
        } else {
            // Check if the signatures correctly sign the transaction head
            if err = txn.Verify(); err != nil { return }
        }

        var txn_input_sum, txn_output_sum uint64

        // Check transaction inputs
        var sender_address foo.U256
        for i, inp := range txn.Inputs {
            if t.GetBillStatus(&inp) != UNSPENT { return errors.New("Input bill is already spent or does not exist") }
            bill := t.GetBill(&inp)
            utils.Assert(bill != nil)
            if ! to_spend.Add(inp) { return errors.New("Two transactions in a block spend the same bill") }
            if i == 0 {
                sender_address = bill.Target
            } else {
                if sender_address != bill.Target { return errors.New("All inputs must be owned by the same person") }
            }
            txn_input_sum += bill.Value
        }

        // Check transaction outputs
        for _, out := range txn.Outputs {
            // TODO: triple-check if the output value does not overflow the counter
            if out.Value == 0 { return errors.New("Bills of value 0 are not allowed") }
            txn_output_sum += out.Value
        }

        if txn_num != 0 {
            if txn_output_sum > txn_input_sum { return errors.New("Transaction output must not be greater then its input") }
        }

        total_input_sum += txn_input_sum
        total_output_sum += txn_output_sum
    }

    // TODO: check for overflows EVERYWHERE
    // The transaction should also generate a reward
    total_input_sum += orchain.GetReward(state.Length)

    // Check if the sums match up
    if total_input_sum != total_output_sum {
        return errors.New("Invalid reward/fees")
    }

    return nil
}
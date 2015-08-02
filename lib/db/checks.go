package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "errors"
    "orwell/lib/utils"
    "github.com/deckarep/golang-set"
)


// TODO: triple-check if the output values do not overflow
func (t *Tx) VerifyTransaction(txn *orchain.Transaction, treat_as_coinbase bool) (input, output uint64, err error) {

    var tid foo.U256
    if tid, err = txn.TryID(); err != nil { return }

    // Check if no other transaction has the same ID
    // Note: https://github.com/bitcoin/bips/blob/master/bip-0030.mediawiki
    if t.GetTransaction(tid) != nil { return 0, 0, errors.New("Transaction ID already in use") }

    if treat_as_coinbase {
        if txn.Proof != nil { return 0, 0, errors.New("The proof is not required/allowed in a coinbase transaction") }
        if len(txn.Inputs) != 0 { return 0, 0, errors.New("Coinbase transaction can have no inputs") }
    } else {
        if err = txn.Verify(); err != nil { return }
    }

    to_spend := mapset.NewSet()

    // Check transaction inputs
    var sender_address foo.U256
    for i, inp := range txn.Inputs {
        if t.GetBillStatus(&inp) != UNSPENT { return 0, 0, errors.New("Input bill is already spent or does not exist") }
        bill := t.GetBill(&inp)
        utils.Assert(bill != nil)
        if ! to_spend.Add(inp) { return 0, 0, errors.New("Bill spent twice") }
        if i == 0 {
            sender_address = bill.Target
        } else {
            if sender_address != bill.Target { return 0, 0, errors.New("All inputs must be owned by the same person") }
        }
        input += bill.Value
    }

    // Check transaction outputs
    for _, out := range txn.Outputs {
        if out.Value == 0 { return 0, 0, errors.New("Bills of value 0 are not allowed") }
        output += out.Value
    }

    if ! treat_as_coinbase {
        if output > input { return 0, 0, errors.New("Transaction output must not be greater then its input") }
    }
    return
}

func (t *Tx) VerifyBlockTransactions(txns []orchain.Transaction) error {

    // We'll collect all inputs and check for duplicates
    to_spend := mapset.NewSet()

    // We also check if there are no transaction duplicates
    txn_ids := mapset.NewSet()

    // Here we'll store the sum of all fees
    var total_input_sum, total_output_sum uint64

    // For each transaction
    for txn_num, txn := range txns {
        txn_input, txn_output, err := t.VerifyTransaction(&txn, txn_num == 0)
        if err != nil { return err }
        total_input_sum += txn_input
        total_output_sum += txn_output
        if ! txn_ids.Add(txn.ID()) { return errors.New("Duplicate transactions in block") }
        for _, inp := range txn.Inputs {
            if ! to_spend.Add(inp) { return errors.New("Two transactions in a block spend the same bill") }
        }
    }

    // TODO: check for overflows EVERYWHERE
    // The transaction should also generate a reward
    state := t.GetState()
    total_input_sum += orchain.GetReward(state.Length)

    // Check if the sums match up
    if total_input_sum != total_output_sum {
        return errors.New("Invalid reward/fees")
    }

    return nil
}

func (t *Tx) VerifyTransactionDoesNotConflict(newTxn *orchain.Transaction, txns []orchain.Transaction) (err error) {
    to_spend := mapset.NewSet()
    for _, txn := range txns {
        for _, inp := range txn.Inputs {
            if ! to_spend.Add(inp) { return errors.New("Bill spent twice") }
        }
    }
    return nil
}

func (t *Tx) VerifyNextBlock(b *orchain.Block) (err error) {
    bid := b.Header.ID()
    state := t.GetState()

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
    // TODO: define the timestamp policy
    if state.Length > 0 {
        previous_header := t.GetHeaderByID(state.Head)
        if b.Header.Timestamp < previous_header.Timestamp {
            return errors.New("Block timestamp is smaller than the previous one")
        }
    }

    if err = b.CheckMerkleRoot(); err != nil { return }
    if err = t.VerifyBlockTransactions(b.Transactions); err != nil { return err }

    return nil
}
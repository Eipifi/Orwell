package db
import (
    "errors"
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/utils"
)

type DBI struct {
    s Storage
}

func NewDB(storage Storage) DB {
    s := &DBI{storage}
    state := s.State()
    if state.Length == 0 {
        utils.Ensure(s.Push(GenesisBlock()))
    }
    return s
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DBI) State() *State {
    return d.s.State()
}

func (d *DBI) GetBlockByID(id foo.U256) *orchain.Block { // TODO: cache
    header := d.GetHeaderByID(id)
    if header == nil { return nil }
    block := &orchain.Block{}
    block.Header = *header
    for _, t := range d.s.GetTransactions(id) {
        txn := d.s.GetTransaction(t)
        utils.Assert(txn != nil)
        block.Transactions = append(block.Transactions, *txn)
    }
    return block
}

func (d *DBI) GetHeaderByID(id foo.U256) *orchain.Header { // TODO: cache
    return d.s.GetHeaderByID(id)
}

func (d *DBI) GetHeaderByNum(num uint64) *orchain.Header {
    return d.s.GetHeaderByNum(num)
}

func (d *DBI) GetNumByID(id foo.U256) *uint64 { // TODO: cache
    return d.s.GetNumByID(id)
}

func (d *DBI) GetIDByNum(num uint64) *foo.U256 {
    return d.s.GetIDByNum(num)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DBI) Push(b *orchain.Block) (err error) {
    // Calculate the block id
    bid := b.Header.ID()
    state := d.State()

    // Check if block is a child of the current head
    if ! foo.Equal(state.Head, b.Header.Previous) {
        return errors.New("The 'Previous' field of the block does not match the stored head")
    }

    // Check if the difficulty value is correct
    if ComputeDifficulty(state.Length, b.Header.Timestamp, d) != b.Header.Difficulty {
        return errors.New("Invalid difficulty value")
    }

    // Check if the block hash meets the specified difficulty
    if ! orchain.HashMeetsDifficulty(bid, b.Header.Difficulty) {
        return errors.New("Block hash does not meet the specified difficulty")
    }

    // Check if the timestamp is correct
    if state.Length > 0 {
        previous_header := d.s.GetHeaderByID(state.Head)
        if b.Header.Timestamp < previous_header.Timestamp {
            return errors.New("Block timestamp is smaller than the previous one")
        }
    }

    // Check if the Merkle root matches (and also if there is at least one transaction)
    if err = b.CheckMerkleRoot(); err != nil { return }

    // Check if the signatures correctly sign the transaction head
    for _, txn := range b.Transactions {
        if err = txn.VerifySignatures(); err != nil { return }
    }

    // Check if input bills are unspent and if the spend proofs are correct
    to_spend := make(map[orchain.BillNumber] bool)
    for _, txn := range b.Transactions {
        for i, inp := range txn.Inputs {
            bill := d.s.GetBill(inp)
            if bill == nil {
                return errors.New("Input bill is already spent or does not exist")
            }
            if _, ok := to_spend[inp]; ok {
                return errors.New("Two transactions in a block spend the same bill")
            }
            to_spend[inp] = true
            pk_id, err := txn.Proofs[i].PublicKey.ID()
            if err != nil { return err }
            if ! foo.Equal(bill.Target, pk_id) {
                return errors.New("The public key does not match the owner of the unspent transaction")
            }
        }
        for _, out := range txn.Outputs {
            if out.Value == 0 { return errors.New("Bills of value 0 are not allowed") }
        }
    }
    var fees uint64 = 0
    // Check if all transactions (except the first) have a legal input/output balance
    // TODO verify if we do not get any uint64 overflows here
    for i := 1; i < len(b.Transactions); i += 1 {
        txn := b.Transactions[i]
        var input_sum uint64 = 0
        var output_sum uint64 = 0
        for _, inp := range txn.Inputs {
            bill := d.s.GetBill(inp)
            utils.Assert(bill != nil) // we already checked if the input bills are unspent, so it should be ok
            input_sum += bill.Value
        }
        for _, out := range txn.Outputs {
            output_sum += out.Value
        }
        if output_sum > input_sum {
            return errors.New("Transaction output sum is bigger than its input sum")
        }
        fees += input_sum - output_sum
    }

    // Calculate the reward for this block number
    var reward uint64 = orchain.GetReward(state.Length)

    // Check if the first transaction correctly grants all the fees (plus reward)
    var txn0_input_sum uint64 = 0
    var txn0_output_sum uint64 = 0
    for _, inp := range b.Transactions[0].Inputs {
        txn0_input_sum += d.s.GetBill(inp).Value
    }
    for _, out := range b.Transactions[0].Outputs {
        txn0_output_sum += out.Value
    }
    if txn0_input_sum + fees + reward != txn0_output_sum {
        return errors.New("Invalid reward/fees")
    }

    // All checks passed, now save the block

    // Insert the block
    utils.Ensure(d.s.PutBlock(b))
    return nil
}

func (d *DBI) Pop() {
    utils.Ensure(d.s.PopBlock())
}
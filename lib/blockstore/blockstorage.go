package blockstore
import (
    "errors"
    "orwell/lib/protocol/orchain"
    "orwell/lib/butils"
)

type BlockStorageImpl struct {
    db Database
}

func NewBlockStore(storage Database) BlockStorage {
    s := &BlockStorageImpl{storage}
    if s.Length() == 0 {
        ensure(s.Push(GenesisBlock()))
    }
    return s
}

func (s *BlockStorageImpl) Head() butils.Uint256 {
    return s.db.Head()
}

func (s *BlockStorageImpl) Length() uint64 {
    return s.db.Length()
}

func (s *BlockStorageImpl) Push(b *orchain.Block) (err error) {
    // Calculate the block id
    bid := b.Header.ID()

    // Check if block is a child of the current head
    if ! butils.Equal(s.Head(), b.Header.Previous) {
        return errors.New("The 'Previous' field of the block dies not match the stored head")
    }

    // Check if the difficulty value is correct
    if s.Length() % orchain.BLOCKS_PER_DIFFICULTY_CHANGE == 0 {
        if s.Length() != 0 {
            // the block needs a recalculated difficulty
            var time_difference = b.Header.Timestamp - s.db.GetHeaderByNum(s.Length() - orchain.BLOCKS_PER_DIFFICULTY_CHANGE).Timestamp
            difficulty_delta := orchain.DifficultyDeltaForTimeDifference(time_difference)
            previous_header := s.db.GetHeaderByID(s.Head())
            if orchain.ApplyDifficultyDelta(previous_header.Difficulty, difficulty_delta) != b.Header.Difficulty {
                return errors.New("The new difficulty is not computed correctly")
            }
        }
    } else {
        previous_header := s.db.GetHeaderByID(s.Head())
        if b.Header.Difficulty != previous_header.Difficulty {
            return errors.New("This block difficulty must be the same as the previous one")
        }
    }

    // Check if the block hash meets the specified difficulty
    if ! orchain.HashMeetsDifficulty(bid, b.Header.Difficulty) {
        return errors.New("Block hash does not meet the specified difficulty")
    }

    // Check if the timestamp is correct
    if s.Length() > 0 {
        previous_header := s.db.GetHeaderByID(s.Head())
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
            bill := s.db.GetBill(inp)
            if bill == nil {
                return errors.New("Input bill is already spent or does not exist")
            }
            if _, ok := to_spend[inp]; ok {
                return errors.New("Two transactions in a block spend the same bill")
            }
            to_spend[inp] = true
            pk_id, err := txn.Proofs[i].PublicKey.ID()
            if err != nil { return err }
            if ! butils.Equal(bill.Target, pk_id) {
                return errors.New("The public key does not match the owner of the unspent transaction")
            }
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
            bill := s.db.GetBill(inp)
            assert(bill != nil) // we already checked if the input bills are unspent, so it should be ok
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
    var reward uint64 = 50 // TODO: implement

    // Check if the first transaction correctly grants all the fees (plus reward)
    var txn0_input_sum uint64 = 0
    var txn0_output_sum uint64 = 0
    for _, inp := range b.Transactions[0].Inputs {
        txn0_input_sum += s.db.GetBill(inp).Value
    }
    for _, out := range b.Transactions[0].Outputs {
        txn0_output_sum += out.Value
    }
    if txn0_input_sum + fees + reward != txn0_output_sum {
        return errors.New("Invalid reward/fees")
    }

    // All checks passed, now save the block

    // Insert the block
    ensure(s.db.PutBlock(b))
    return nil
}

func (s *BlockStorageImpl) Pop() {
    ensure(s.db.PopBlock())
}

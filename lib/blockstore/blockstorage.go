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
    head, _ := s.db.FetchHead()
    return head
}

func (s *BlockStorageImpl) Length() uint64 {
    _, l := s.db.FetchHead()
    return l
}

func (s *BlockStorageImpl) Push(b *orchain.Block) (err error) {
    // Calculate the block id
    bid := b.Header.ID()

    // Check if block is a child of the current head
    if butils.Compare(s.Head(), b.Header.Previous) != 0 {
        return errors.New("The 'Previous' field of the block dies not match the stored head")
    }

    // Check if the block hash meets the specified difficulty
    if ! orchain.HashMeetsDifficulty(bid, b.Header.Difficulty) {
        return errors.New("Block hash does not meet the specified difficulty")
    }

    // Check if the difficulty value is correct
    // TODO: Verify difficulty

    // Check if the timestamp is correct
    // TODO: Verify timestamp

    // Check if the Merkle root matches (and also if there is at least one transaction)
    if err = b.CheckMerkleRoot(); err != nil { return }

    // Check if the signatures correctly sign the transaction head
    for _, txn := range b.Transactions {
        if err = txn.VerifySignatures(); err != nil { return }
    }

    // Check if input bills are unspent and if the spend proofs are correct
    for _, txn := range b.Transactions {
        for i, inp := range txn.Inputs {
            bill := s.db.FetchUnspentBill(inp)
            if bill == nil {
                return errors.New("Input bill is already spent or does not exist")
            }
            pk_id, err := txn.Proofs[i].PublicKey.ID()
            if err != nil { return err }
            if butils.Compare(bill.Target, pk_id) != 0 {
                return errors.New("The public key does not match the owner of the unspent transaction")
            }
        }
    }

    // Check if no transaction output is spent twice in the block transactions
    // TODO

    var fees uint64 = 0
    // Check if all transactions (except the first) have a legal input/output balance
    // TODO verify if we do not get any uint64 overflows here
    for i := 1; i < len(b.Transactions); i += 1 {
        txn := b.Transactions[i]
        var input_sum uint64 = 0
        var output_sum uint64 = 0
        for _, inp := range txn.Inputs {
            bill := s.db.FetchUnspentBill(inp)
            assert(bill != nil) // we already checked if the inputs are unspent, so it should be ok
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
        txn0_input_sum += s.db.FetchUnspentBill(inp).Value
    }
    for _, out := range b.Transactions[0].Outputs {
        txn0_output_sum += out.Value
    }
    if txn0_input_sum + fees + reward != txn0_output_sum {
        return errors.New("Invalid reward/fees")
    }

    // All checks passed, now save the block

    // Insert the block
    ensure(s.db.StoreHeader(&b.Header))

    // Update the head
    s.db.StoreHead(bid, s.Length() + 1)

    // Insert the transactions
    for _, txn := range b.Transactions {
        tid, _ := txn.ID()
        ensure(s.db.StoreTransaction(&txn))
        for _, inp := range txn.Inputs {
            s.db.SpendBill(inp)
        }
        for i, out := range txn.Outputs {
            s.db.StoreUnspentBill(orchain.BillNumber{tid, uint64(i)}, out)
        }
    }

    return nil
}

func (s *BlockStorageImpl) Pop() {
    if s.Length() <= 1 { return } // thou shalt not remove the genesis block


}

package miner
import (
    "orwell/lib/blockstore"
    "orwell/lib/protocol/orchain"
    "orwell/lib/timing"
    "orwell/lib/butils"
    "log"
    "time"
    "orwell/lib/logging"
)

// Note:
// Txn buffer stores transactions that the miner will try to sign.
// Every time a new block arrives (external or internal)

// This class is responsible for keeping the miner busy.
// It watches the results, and updates the objective block each second.
type MiningSupervisor struct {
    log *log.Logger
    miner Miner
    account butils.Uint256
    storage blockstore.BlockStorage
    txn_buffer []orchain.Transaction // must be refreshed every time a head changes
}

func NewSupervisor(account butils.Uint256, storage blockstore.BlockStorage) *MiningSupervisor {
    s := &MiningSupervisor{}
    s.log = logging.GetLogger("")
    s.account = account
    s.storage = storage
    return s
}

func (s *MiningSupervisor) createObjective() (b orchain.Block) {
    num := s.storage.Length()
    b.Header.Previous = s.storage.Head()
    b.Header.Timestamp = timestamps.CurrentTimestamp()
    b.Header.Difficulty = blockstore.ComputeDifficulty(num, b.Header.Timestamp, s.storage)
    b.Transactions = []orchain.Transaction{
        orchain.Transaction{
            Outputs: []orchain.Bill {
                orchain.Bill{
                    Target: s.account,
                    Value: 50,
                },
            },
        },
    }
    b.Transactions = append(b.Transactions, s.txn_buffer...)
    b.ComputeMerkleRoot()
    return
}

func (s *MiningSupervisor) Run() {
    s.miner.SetObjective(s.createObjective())
    s.miner.RunWorkers(1)
    for {
        select {
            case block := <- s.miner.Results:
                err := s.storage.Push(&block)
                if err == nil {
                    s.log.Printf("Successfully mined and stored the block %v", block.Header.ID())
                } else {
                    s.log.Printf("Failed to store the mined block - reason: %v", err)
                }
            case <-time.After(time.Second * 1):
                // block was not mined, just update the objective.
        }
        s.miner.SetObjective(s.createObjective())
    }
}

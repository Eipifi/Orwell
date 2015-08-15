package miner
import (
    "orwell/lib/db"
    "orwell/lib/protocol/orchain"
    "orwell/lib/timing"
    "orwell/lib/foo"
    "math/rand"
    "log"
    "orwell/orchain/serv"
    "orwell/lib/utils"
)

type SimpleMiner struct {
    run bool
}

func StartMiner(wallet foo.U256) (*SimpleMiner) {
    m := &SimpleMiner{true}
    go m.work(wallet)
    return m
}

func (m *SimpleMiner) work(wallet foo.U256) {
    for m.run {
        block := prepareBlock(wallet)
        if trySign(block, 1000000) {
            err := serv.SyncMgr().PushBlock(block)
            if err != nil {
                log.Printf("Mined block was not saved: %v", err)
            } else {
                log.Printf("Mined block %v", block.Header.ID())
            }
        }
    }
}

func (m *SimpleMiner) Stop() {
    m.run = false
}

func trySign(block *orchain.Block, iterations int) bool {
    block.Header.Nonce = uint64(rand.Uint32()) << 32
    for i := iterations; i > 0; i -= 1 {
        id := block.Header.ID()
        if orchain.HashMeetsDifficulty(id, block.Header.Difficulty) {
            return true
        }
        block.Header.Nonce += 1
    }
    return false
}

func prepareBlock(wallet foo.U256) (block *orchain.Block) {
    block = &orchain.Block{}

    db.Get().View(func(t *db.Tx) {
        state := t.GetState()
        block.Header.Previous = state.Head
        block.Header.Timestamp = timestamps.CurrentTimestamp()
        block.Header.Difficulty = t.GetDifficulty(state.Length)
        block.Transactions = []orchain.Transaction{
            orchain.Transaction{
                Outputs: []orchain.Bill {
                    orchain.Bill{
                        Target: wallet,
                        Value: orchain.GetReward(state.Length),
                    },
                },
                Payload: orchain.PayloadLabelString("Block #" + string(state.Length)),
            },
        }
        block.Transactions = append(block.Transactions, t.UnconfirmedTransactions()...)

        for i := 1; i < len(block.Transactions); i += 1 {
            fee, err := t.ComputeTransactionFee(&block.Transactions[i])
            utils.Ensure(err)
            block.Transactions[0].Outputs[0].Value += fee
        }

        block.Domains = t.DomainsToRegister(block.Transactions)
        block.ComputeMerkleRoot()
    })

    return
}
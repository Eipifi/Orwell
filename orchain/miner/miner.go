package miner
import (
    "orwell/lib/protocol/orchain"
    "time"
    "math/rand"
    "orwell/lib/foo"
    "orwell/lib/timing"
    "log"
    "orwell/lib/logging"
    "orwell/lib/db"
)

const MINER_ITERATIONS = 500000
const MINER_WORKERS = 1

type Miner struct {
    log *log.Logger
    target foo.U256
    objective *orchain.Block
    run bool
}

func NewMiner(target foo.U256) *Miner {
    m := &Miner{}
    m.log = logging.GetStdLogger("")
    m.target = target
    m.Run(false)
    m.update()
    for i := 0; i < MINER_WORKERS; i += 1 {
        go m.work()
    }
    go m.apply()
    return m
}

func (m *Miner) Run(run bool) {
    m.run = run
}

func (m *Miner) work() {
    for {
        if ! m.run {
            time.Sleep(time.Second)
            continue
        }

        block := *(m.objective)
        block.Header.Nonce = uint64(rand.Uint32()) << 32
        for i := 0; i < MINER_ITERATIONS; i += 1 {
            id := block.Header.ID()
            if orchain.HashMeetsDifficulty(id, block.Header.Difficulty) {
                m.deliver(block)
                break
            }
            block.Header.Nonce += 1
        }
    }
}

func (m *Miner) deliver(block orchain.Block) {
    err := db.Get().Push(&block)
    if err == nil {
        log.Printf("Mined block id=%v df=%v \n", block.Header.ID(), block.Header.Difficulty)
    } else {
        log.Printf("Error while applying block: %v", err)
    }
    m.update()
}

func (m *Miner) apply() {
    for {
        m.update()
        time.Sleep(time.Second)
    }
}

func (m *Miner) update() {
    state := db.Get().State()
    num := state.Length
    block := orchain.Block{}
    block.Header.Previous = state.Head
    block.Header.Timestamp = timestamps.CurrentTimestamp()
    block.Header.Difficulty = db.Get().Difficulty()
    block.Transactions = []orchain.Transaction{
        orchain.Transaction{
            Outputs: []orchain.Bill {
                orchain.Bill{
                    Target: m.target,
                    Value: orchain.GetReward(num),
                },
            },
            Label: "Block #" + string(num),
        },
    }
    block.ComputeMerkleRoot()
    m.objective = &block
}
package miner
import (
    "orwell/lib/protocol/orchain"
    "math/rand"
)

// The main problem is preserving both the header AND the transaction array.
// The transaction array is a slice
// we only append new transactions, so each slice will remain the same
// So we only need to

const ATTEMPTS_PER_UNIT_WORK = 500000

type Miner struct {
    objective orchain.Block
    Results chan orchain.Block
    running bool
}

func (m *Miner) SetObjective(block orchain.Block) {
    // todo: think about atomic worker reads during block update
    m.objective = block
}

func (m *Miner) RunWorkers(n int) {
    m.Results = make(chan orchain.Block)
    m.running = true
    for i := 0; i < n; i += 1 {
        go m.worker()
    }
}

func (m *Miner) StopWorkers() {
    m.running = false
}

func (m *Miner) worker() {
    for (m.running) {
        m.work()
    }
}

func (m *Miner) work() {
    // Make a local copy of the block (transactions will be shared by slice, but that's not a problem)
    block := m.objective
    // Randomize a nonce prefix (we're changing only the last 32 bits)
    block.Header.Nonce = uint64(rand.Uint32()) << 32
    // Compute the target
    target := orchain.DifficultyToTarget(block.Header.Difficulty)
    // We attempt a million hashes (about one second right now, maybe try some dynamic performance tracking)
    for i := 0; i < 1000000; i += 1 {
        // this is a clean implementation, not an efficient one
        id := block.Header.ID()
        if orchain.HashMeetsTarget(id, target) {
            m.Results <- block
            return
        }
        block.Header.Nonce += 1
    }
}


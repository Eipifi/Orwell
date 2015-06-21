package miner
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/butils"
    "crypto/sha256"
)

func Mine(block *orchain.Block) {
    // This is not an efficient implementation, nor it should be
    target := orchain.DifficultyToTarget(block.Difficulty)
    for {
        buf, _ := butils.WriteToBytes(block)
        id := sha256.Sum256(buf)
        if orchain.HashMeetsTarget(id, target) { break }
        block.Nonce += 1
    }
}

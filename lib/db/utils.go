package db
import (
    "orwell/lib/foo"
    "orwell/lib/protocol/orchain"
)

func ComputeDifficulty(num uint64, ts uint64, d DB) foo.U256 {
    if num == 0 { return GENESIS_DIFFICULTY }
    if num % orchain.BLOCKS_PER_DIFFICULTY_CHANGE == 0 {
        prev := d.GetHeaderByNum(num - orchain.BLOCKS_PER_DIFFICULTY_CHANGE)
        return orchain.UpdateDifficulty(prev.Difficulty, ts - prev.Timestamp)
    } else {
        return d.GetHeaderByNum(num - 1).Difficulty
    }
}
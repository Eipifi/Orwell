package db
import (
    "orwell/lib/foo"
    "orwell/lib/protocol/orchain"
    "orwell/lib/utils"
)

func (t *Tx) GetDifficulty(num uint64) foo.U256 {
    if num <= 1 { return orchain.GENESIS_DIFFICULTY }
    prev := t.GetHeaderByNum(num - 1)
    utils.Assert(prev != nil)
    if num % orchain.BLOCKS_PER_DIFFICULTY_CHANGE == 1 {
        prev2 := t.GetHeaderByNum(num - 1 - orchain.BLOCKS_PER_DIFFICULTY_CHANGE)
        return orchain.UpdateDifficulty(prev.Difficulty, prev.Timestamp - prev2.Timestamp)
    } else {
        return prev.Difficulty
    }
}
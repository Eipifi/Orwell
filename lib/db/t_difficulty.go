package db
import (
    "orwell/lib/foo"
    "orwell/lib/protocol/orchain"
)

func (t *Tx) GetDifficulty() foo.U256 {
    state := t.GetState()
    num := state.Length
    if num <= 1 { return orchain.GENESIS_DIFFICULTY }
    head := t.GetHeaderByID(state.Head)
    if num % orchain.BLOCKS_PER_DIFFICULTY_CHANGE == 1 {
        prev := t.GetHeaderByNum(num - 1 - orchain.BLOCKS_PER_DIFFICULTY_CHANGE)
        return orchain.UpdateDifficulty(head.Difficulty, head.Timestamp - prev.Timestamp)
    } else {
        return head.Difficulty
    }
}
package db
import (
    "orwell/lib/foo"
    "orwell/lib/protocol/orchain"
    "github.com/boltdb/bolt"
)

func Difficulty(t *bolt.Tx) foo.U256 {
    state := GetState(t)
    num := state.Length
    if num <= 1 { return GENESIS_DIFFICULTY }
    head := GetHeaderByID(t, state.Head)
    if num % orchain.BLOCKS_PER_DIFFICULTY_CHANGE == 1 {
        prev := GetHeaderByNum(t, num - 1 - orchain.BLOCKS_PER_DIFFICULTY_CHANGE)
        return orchain.UpdateDifficulty(head.Difficulty, head.Timestamp - prev.Timestamp)
    } else {
        return head.Difficulty
    }
}
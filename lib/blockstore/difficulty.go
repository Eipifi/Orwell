package blockstore
import "orwell/lib/protocol/orchain"

func ComputeDifficulty(num uint64, ts uint64, s BlockStorage) uint16 {
    if num == 0 { return GENESIS_DIFFICULTY }
    if num % orchain.BLOCKS_PER_DIFFICULTY_CHANGE == 0 {
        prev := s.GetHeaderByNum(num - orchain.BLOCKS_PER_DIFFICULTY_CHANGE)
        time_difference := ts - prev.Timestamp
        difficulty_delta := orchain.DifficultyDeltaForTimeDifference(time_difference)
        return orchain.ApplyDifficultyDelta(prev.Difficulty, difficulty_delta)
    } else {
        return s.GetHeaderByNum(num - 1).Difficulty
    }
}
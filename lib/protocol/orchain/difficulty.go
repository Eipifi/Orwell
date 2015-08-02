package orchain
import "orwell/lib/foo"

const BLOCKS_PER_DIFFICULTY_CHANGE = 32
const SECONDS_PER_BLOCK = 10
const DIFFICULTY_FACTOR = 4

func DifficultyToTarget(difficulty foo.U256) foo.U256 {
    if difficulty == foo.ZERO || difficulty == foo.ONE {
        return foo.MAX
    }
    return difficulty.Invert256()
}

func HashMeetsTarget(hash, target foo.U256) bool {
    return foo.Compare(hash, target) < 0
}

func HashMeetsDifficulty(hash foo.U256, difficulty foo.U256) bool {
    return HashMeetsTarget(hash, DifficultyToTarget(difficulty))
}

func UpdateDifficulty(difficulty foo.U256, delta_obtained uint64) foo.U256 {
    var delta_expected uint64 = BLOCKS_PER_DIFFICULTY_CHANGE * SECONDS_PER_BLOCK

    if delta_obtained > DIFFICULTY_FACTOR * delta_expected {
        delta_obtained = DIFFICULTY_FACTOR * delta_expected
    }

    if delta_obtained < delta_expected / DIFFICULTY_FACTOR {
        delta_obtained = delta_expected / DIFFICULTY_FACTOR
    }

    result := difficulty.MulDiv64(delta_expected, delta_obtained)
    if result == foo.ZERO {
        return foo.ONE
    } else {
        return result
    }
}
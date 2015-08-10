package db
import (
    "orwell/lib/protocol/orchain"
    "errors"
    "orwell/lib/foo"
)

func CheckHeaderBasics(t *Tx, s *State, h *orchain.Header) (err error) {
    // Check if no other block has the same id
    // Note: while the chance of this happening is astronomically low, we still check this.
    // SHA256 might get broken at some point in the future.
    if t.GetHeaderByID(h.ID()) != nil {
        return errors.New("Block with this ID already exists")
    }

    // Check if block is a child of the current head
    if ! foo.Equal(s.Head, h.Previous) {
        return errors.New("The 'Previous' field of the block does not match the stored head")
    }
    return
}

func CheckHeaderDifficulty(t *Tx, s *State, h *orchain.Header) (err error) {
    // Check if the difficulty value is correct
    if t.GetDifficulty(s.Length) != h.Difficulty {
        return errors.New("Invalid difficulty value")
    }

    // Check if the block hash meets the specified difficulty
    if ! orchain.HashMeetsDifficulty(h.ID(), h.Difficulty) {
        return errors.New("Block hash does not meet the specified difficulty")
    }
    return
}

func CheckHeaderTimestamp(t *Tx, s *State, h *orchain.Header) (err error) {
    // Check if the timestamp is correct
    // TODO: define the timestamp policy
    if s.Length > 0 {
        previous_header := t.GetHeaderByID(s.Head)
        if h.Timestamp < previous_header.Timestamp {
            return errors.New("Block timestamp is smaller than the previous one")
        }
    }
    return
}
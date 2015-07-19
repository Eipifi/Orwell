package merkle
import (
    "orwell/lib/crypto/hash"
    "orwell/lib/foo"
)

func Compute(IDs []foo.U256) foo.U256 {

    l := len(IDs)

    // Assert the slice is not empty
    if l == 0 {
        panic("Cannot compute a Merkle tree of 0 elements")
    }

    // If there is only one element, return it
    if l == 1 {
        return IDs[0]
    }

    // If the number of elements is odd, copy the last one
    if l % 2 == 1 {
        IDs = append(IDs, IDs[l - 1])
        l += 1
    }

    // Create a new slice where the hashes will be held
    Sums := make([]foo.U256, l/2, l/2)

    // Combine pairs of hashes
    for i := 0; i < l/2; i += 1 {
        Sums[i] = combine(IDs[2*i], IDs[2*i + 1])
    }

    return Compute(Sums)
}

func combine(a, b foo.U256) foo.U256 {
    buf := append(a[:], b[:]...)
    return hash.Hash(buf)
}
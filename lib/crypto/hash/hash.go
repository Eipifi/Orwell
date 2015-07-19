package hash
import (
    "crypto/sha256"
    "orwell/lib/butils"
    "orwell/lib/foo"
)

func Hash(data []byte) foo.U256 {
    return sha256.Sum256(data)
}

func HashBytes(data []byte) []byte {
    id := Hash(data)
    return id[:]
}

func HashOf(w butils.Writable) (id foo.U256, err error) {
    buf, err := butils.WriteToBytes(w)
    if err != nil { return }
    return Hash(buf), nil
}
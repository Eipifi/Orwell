package sig

import (
	"crypto/sha256"
	"encoding/hex"
)

const HashLengthBytes = 32
const HashLengthBits = HashLengthBytes * 8
type ID [HashLengthBytes]byte

func Hash(data []byte) *ID {
	var id ID = sha256.Sum256(data)
	return &id
}

func HashSlice(data []byte) []byte {
	h := Hash(data)
	return h[:]
}

func (id *ID) String() string {
	return hex.EncodeToString(id[:])
}
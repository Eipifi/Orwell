package sig

import "crypto/sha256"

const HashLengthBytes = 32
const HashLengthBits = HashLengthBytes * 8
type ID [HashLengthBytes]byte

func Hash(data []byte) ID {
	return sha256.Sum256(data)
}

func HashSlice(data []byte) []byte {
	h := Hash(data)
	return h[:]
}

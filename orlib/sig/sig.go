// Package sig handles signatures and public keys
package sig
import "orwell/orlib/protocol/types"

type PubKey interface {
    Serializable
    types.IDer
    Verify(data []byte, signature []byte) bool
}

type PrvKey interface {
    Serializable
    PublicPart() PubKey
    Sign(data []byte) []byte
}
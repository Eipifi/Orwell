// Package sig handles signatures and public keys
package sig

type PubKey interface {
    Serializable
    Id() *ID
    Verify(data []byte, signature []byte) bool
}

type PrvKey interface {
    Serializable
    PublicPart() PubKey
    Sign(data []byte) []byte
}
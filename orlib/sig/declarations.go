package sig

type Key interface {
	Id() ID
	Marshal() []byte
}

type PubKey interface {
	Key
	Verify(data []byte, signature []byte) bool
}

type PrvKey interface {
	Key
	PublicPart() PubKey
	Sign(data []byte) []byte
}

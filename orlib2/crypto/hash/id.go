package hash
import (
    "io"
    "orwell/orlib2/butils"
    "crypto/sha256"
    "encoding/hex"
)

const (
    ByteLength = 32
    BitLength = ByteLength * 8
)

type ID [ByteLength]byte

func (i *ID) Read(r io.Reader) error {
    return butils.ReadFull(r, i[:])
}

func (i *ID) Write(w io.Writer) error {
    return butils.WriteFull(w, i[:])
}

func NewId(data []byte) *ID {
    var id ID = sha256.Sum256(data)
    return &id
}

func Hash(data []byte) []byte {
    id := NewId(data)
    return id[:]
}

func (id *ID) String() string {
    return hex.EncodeToString(id[:])
}
package hash
import (
    "io"
    "orwell/orlib/butils"
    "crypto/sha256"
    "encoding/hex"
    "errors"
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

func HexToID(h string) (*ID, error) {
    b, err := hex.DecodeString(h)
    if err != nil { return nil, err }
    if len(b) != ByteLength { return nil, errors.New("Invalid hex length") }
    id := ID{}
    copy(id[:], b)
    return &id, nil
}

func Equal(id1 *ID, id2 *ID) bool {
    panic("Not implemented")
}
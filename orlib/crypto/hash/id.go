package hash
import (
    "io"
    "orwell/orlib/butils"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "bytes"
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

func New(data []byte) ID {
    return sha256.Sum256(data)
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

func Compare(a, b ID) int {
    return bytes.Compare(a[:], b[:])
}

func LeftCloser(a, b, c ID) bool {
    for i := 0; i < ByteLength; i++ {
        dl, dr := dist(a[i], b[i]), dist(b[i], c[i])
        if dl < dr { return true }
        if dl > dr { return false }
    }
    return true
}

func dist(a, b byte) int {
    diff := int(a) - int(b)
    if diff < 0 { diff = -diff }
    ndiff := 256 - diff
    if diff < ndiff { return diff }
    return ndiff
}
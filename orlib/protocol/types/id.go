package types
import (
    "orwell/orlib/comm"
    "crypto/sha256"
    "encoding/hex"
    "errors"
)

const HashLengthBytes = 32
const HashLengthBits = HashLengthBytes * 8
type ID [HashLengthBytes]byte

type IDer interface {
    Id() *ID
}

func (i *ID) Read(r *comm.Reader) error {
    return r.ReadTo(i[:])
}

func (i *ID) Write(w *comm.Writer) {
    w.Write(i[:])
}

func Hash(data []byte) *ID {
    var id ID = sha256.Sum256(data)
    return &id
}

func HashSlice(data []byte) []byte {
    h := Hash(data)
    return h[:]
}

func HexToID(h string) (*ID, error) {
    b, err := hex.DecodeString(h)
    if err != nil { return nil, err }
    if len(b) != HashLengthBytes { return nil, errors.New("Invalid hex length") }
    id := ID{}
    copy(id[:], b)
    return &id, nil
}

func (id *ID) String() string {
    return hex.EncodeToString(id[:])
}
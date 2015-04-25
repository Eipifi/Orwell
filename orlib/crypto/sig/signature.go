package sig
import (
    "math/big"
    "io"
    "github.com/eipifi/asn1"
    "errors"
)

type Signature struct {
    R, S *big.Int
}

func (s *Signature) Read(r io.Reader) error {
    return asn1.UnmarshalFromReader(s, r)
}

func (s *Signature) Write(w io.Writer) error {
    return asn1.MarshalToWriter(s, w)
}

func (s *Signature) ReadBytes(data []byte) error {
    rest, err := asn1.Unmarshal(data, s)
    if err != nil { return err }
    if len(rest) > 0 { return errors.New("Unread bytes remaining") }
    return nil
}

func (s *Signature) WriteBytes() ([]byte, error) {
    return asn1.Marshal(s)
}


package sig
import (
    "math/big"
    "errors"
    "encoding/asn1"
)

type Signature struct {
    R, S *big.Int
}

func (s *Signature) ReadBytes(data []byte) error {
    rest, err := asn1.Unmarshal(data, s)
    if err != nil { return err }
    if len(rest) > 0 { return errors.New("Unread bytes remaining") }
    return nil
}

func (s *Signature) WriteBytes() ([]byte, error) {
    return asn1.Marshal(*s)
}


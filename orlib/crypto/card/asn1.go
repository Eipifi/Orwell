package card
import (
    "errors"
    "encoding/asn1"
)

type asn1Card struct {
    Key []byte
    Payload Payload
    Signature []byte
}

func (c *asn1Card) ReadBytes(data []byte) error {
    return unmarshalAll(data, c)
}

func (c *asn1Card) WriteBytes() ([]byte, error) {
    return asn1.Marshal(*c)
}

func unmarshalAll(data []byte, p interface{}) error {
    rest, err := asn1.Unmarshal(data, p)
    if err != nil { return err }
    if len(rest) > 0 { return errors.New("Bytes remaining") }
    return nil
}
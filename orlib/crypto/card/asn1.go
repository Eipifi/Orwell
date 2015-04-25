package card
import (
    "io"
    "github.com/eipifi/asn1"
)

type asn1Card struct {
    Key []byte
    Payload Payload
    Signature []byte
}

func (c *asn1Card) Read(r io.Reader) error {
    return asn1.UnmarshalFromReader(c, r)
}

func (c *asn1Card) Write(w io.Writer) error {
    return asn1.MarshalToWriter(c, w)
}

func (c *asn1Card) ReadBytes(data []byte) error {
    return asn1.Unmarshal(data, c)
}

func (c *asn1Card) WriteBytes() ([]byte, error) {
    return asn1.Marshal(c)
}
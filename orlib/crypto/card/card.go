package card
import (
    "orwell/orlib/crypto/sig"
    "io"
    "errors"
)

type Card struct {
    Key *sig.PublicKey
    Payload *Payload
    Signature *sig.Signature
}

func (c *Card) Read(r io.Reader) (err error) {
    ac := &asn1Card{}
    if err = ac.Read(r); err != nil { return }
    return c.readAsn1Card(ac)
}

func (c *Card) Write(w io.Writer) error {
    ac, err := c.writeAsn1Card();
    if err != nil { return err }
    return ac.Write(w)
}

func (c *Card) ReadBytes(data []byte) (err error) {
    ac := &asn1Card{}
    if err = ac.ReadBytes(data); err != nil { return }
    return c.readAsn1Card(ac)
}

func (c *Card) WriteBytes() ([]byte, error) {
    ac, err := c.writeAsn1Card();
    if err != nil { return err }
    return ac.WriteBytes()
}

func (c *Card) readAsn1Card(ac *asn1Card) (err error) {
    c.Key = &sig.PublicKey{}
    if err = c.Key.ReadBytes(ac.Key); err != nil { return }
    c.Signature = &sig.Signature{}
    if err = c.Signature.ReadBytes(ac.Signature); err != nil { return }
    return c.Verify()
}

func (c *Card) writeAsn1Card() (ac *asn1Card, err error) {
    ac = &asn1Card{}
    ac.Payload = c.Payload
    if ac.Key, err = c.Key.WriteBytes(); err != nil { return }
    if ac.Signature, err = c.Signature.WriteBytes(); err != nil { return }
    return
}

func (c *Card) Verify() error {
    return errors.New("Not yet implemented")
}
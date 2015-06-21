package card
import (
    "errors"
    "orwell/lib/crypto/sig"
)

type Card struct {
    Key *sig.PubKey
    Payload *Payload
    Signature *sig.Signature
}

func (c *Card) ReadBytes(data []byte) (err error) {
    ac := &asn1Card{}
    if err = ac.ReadBytes(data); err != nil { return }
    return c.readAsn1Card(ac)
}

func (c *Card) WriteBytes() ([]byte, error) {
    ac, err := c.writeAsn1Card();
    if err != nil { return nil, err }
    return ac.WriteBytes()
}

func (c *Card) readAsn1Card(ac *asn1Card) (err error) {
    c.Key = &sig.PubKey{}
    if err = c.Key.ReadBytes(ac.Key); err != nil { return }
    c.Signature = &sig.Signature{}
    if err = c.Signature.ReadBytes(ac.Signature); err != nil { return }
    c.Payload = &ac.Payload
    return c.Verify()
}

func (c *Card) writeAsn1Card() (ac *asn1Card, err error) {
    ac = &asn1Card{}
    ac.Payload = *(c.Payload)
    if ac.Key, err = c.Key.WriteBytes(); err != nil { return }
    if ac.Signature, err = c.Signature.WriteBytes(); err != nil { return }
    return
}

func (c *Card) Verify() error {
    if c.Key == nil { return errors.New("Key not set") }
    if c.Payload == nil { return errors.New("Payload not set") }
    if c.Signature == nil { return errors.New("Signature not set") }
    buf, err := c.Payload.WriteBytes()
    if err != nil { return err }
    return c.Key.Verify(buf, c.Signature)
}
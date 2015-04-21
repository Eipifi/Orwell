package card
import (
    "time"
    "orwell/orlib/sig"
    "errors"
)

func (c *Card) ExpirationDate() time.Time {
    return time.Unix(c.Payload.Expires, 0)
}

func (c *Card) PubKey() sig.PubKey {
    if pub, err := sig.ParsePubKey(c.Key); err == nil { return pub }
    return nil
}

func (c *Card) Sign(prv sig.PrvKey) (err error) {
    var buf []byte
    if buf, err = c.Payload.MarshalBinary(); err != nil { return }
    c.Signature = prv.Sign(buf)
    c.Key = prv.PublicPart().Serialize()
    return
}

func (c *Card) Verify() (err error) {
    var buf []byte
    if buf, err = c.Payload.MarshalBinary(); err != nil { return }
    pub := c.PubKey()
    if pub == nil { return errors.New("Invalid public key") }
    if pub.Verify(buf, c.Signature) { return nil }
    return errors.New("Verification failed")
}

func UnmarshalOnlyJSON(data []byte) (c *Card, err error) {
    c = &Card{}
    return c, c.Payload.UnmarshalJSON(data)
}
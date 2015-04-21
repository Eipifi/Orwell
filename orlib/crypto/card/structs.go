package card

import (
    "encoding/asn1"
    "encoding/json"
    "orwell/orlib/crypto/armor"
)

type Card struct {
    Key []byte
    Payload Payload
    Signature []byte
}

type Payload struct {
    Version int64 `json:"version"`
    Expires int64 `json:"expires"`
    Records []Record `json:"records"`
}

type Record struct {
    Key string `asn1:"utf8" json:"key"`
    Type string `asn1:"utf8" json:"type"`
    Value string `asn1:"utf8" json:"value"`
}

func (c *Card) MarshalBinary() ([]byte, error) {
    return asn1.Marshal(c)
}

func (c *Card) UnmarshalBinary(b []byte) (err error) {
    if err = UnmarshalAll(b, c); err != nil { return }
    return c.Verify()
}

func (p *Payload) MarshalBinary() ([]byte, error) {
    return asn1.Marshal(p)
}

func (p *Payload) UnmarshalBinary(b []byte) error {
    return UnmarshalAll(b, p)
}

func (p *Payload) MarshalJSON() ([]byte, error) {
    return json.MarshalIndent(p, "", "    ") // 4 spaces
}

func (p *Payload) UnmarshalJSON(b []byte) error {
    return json.Unmarshal(b, p)
}

func (c *Card) UnmarshalPEM(b []byte) error {
    return armor.Unmarshal(b, c)
}

func (c *Card) MarshalPEM() ([]byte, error) {
    return armor.Marshal("CARD", c)
}
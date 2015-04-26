package card
import (
    "encoding/json"
    "orwell/orlib/crypto/sig"
    "encoding/asn1"
)

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

func (p *Payload) ReadBytes(data []byte) error {
    return unmarshalAll(data, p)
}

func (p *Payload) WriteBytes() ([]byte, error) {
    return asn1.Marshal(*p)
}

func (p *Payload) ReadJSON(data []byte) error {
    return json.Unmarshal(data, p)
}

func (p *Payload) WriteJSON() ([]byte, error) {
    return json.MarshalIndent(p, "", "    ")
}

func (p *Payload) Sign(key *sig.PrivateKey) (c *Card, err error) {
    c = &Card{}
    c.Payload = p
    c.Key = key.PublicPart()
    c.Signature, err = key.SignByteWritable(p)
    return
}
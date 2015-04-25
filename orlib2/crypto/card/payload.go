package card
import (
    "io"
    "github.com/eipifi/asn1"
    "encoding/json"
    "orwell/orlib2/crypto/sig"
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

func (p *Payload) Read(r io.Reader) error {
    return asn1.UnmarshalFromReader(p, r)
}

func (p *Payload) Write(w io.Writer) error {
    return asn1.MarshalToWriter(p, w)
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
    c.Signature, err = key.SignWritable(p)
    return
}
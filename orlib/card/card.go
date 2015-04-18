package card

import (
    "time"
    "orwell/orlib/sig"
    "encoding/asn1"
    "encoding/json"
    "errors"
)

// The card structure
type Card struct {
    Key sig.PubKey
    Payload *Payload
    Signature []byte
}

// ASN1 compatible card structure
type Asn1Card struct {
    Key []byte
    Payload Payload
    Signature []byte
}

// Card payload (the part signed with private key)
type Payload struct {
    Version int64 `json:"version"`
    Expires int64 `json:"expires"`
    Records []Record `json:"records"`
}

// Record type
type Record struct {
    Key string `asn1:"utf8" json:"key"`
    Type string `asn1:"utf8" json:"type"`
    Value string `asn1:"utf8" json:"value"`
}

////////////////////////////////////

func (card *Card) Marshal() ([]byte, error) {
    ac := Asn1Card{card.Key.Serialize(), *(card.Payload), card.Signature}
    return asn1.Marshal(ac)
}

func (card *Card) Sign(key sig.PrvKey) {
    card.Signature = key.Sign(card.Payload.MarshalASN1())
    card.Key = key.PublicPart()
}

func (card *Card) Verify() bool {
    return card.Key.Verify(card.Payload.MarshalASN1(), card.Signature)
}

func Unmarshal(data []byte) (*Card, error) {
    ac := Asn1Card{}
    rest, err := asn1.Unmarshal(data, &ac)
    if len(rest) > 0 {
        return nil, errors.New("Unnecesary bytes remaining")
    }
    if err != nil {
        return nil, err
    }
    return UnmarshalStruct(&ac)
}

func UnmarshalStruct(ac *Asn1Card) (*Card, error) {
    key, err := sig.ParsePubKey(ac.Key)
    if err != nil {
        return nil, err
    }
    card := Card{key, &(ac.Payload), ac.Signature}
    if !card.Verify() {
        return nil, errors.New("Card verification failed")
    }
    return &card, nil
}

func (card *Card) ExpirationDate() time.Time {
    return time.Unix(card.Payload.Expires, 0)
}

func (payload *Payload) MarshalASN1() []byte {
    b, err := asn1.Marshal(*payload)
    if err != nil {
        panic(err)
    }
    return b
}

func (payload *Payload) MarshalJSON() []byte {
    b, err := json.MarshalIndent(*payload, "", "    ")
    if err != nil {
        panic(err)
    }
    return b
}

func Create(jsonPayload []byte, key sig.PrvKey) (*Card, error) {
    jp := Payload{}
    if err := json.Unmarshal(jsonPayload, &jp); err != nil {
        return nil, err
    }
    card := Card{nil, &jp, nil}
    card.Sign(key)
    return &card, nil
}
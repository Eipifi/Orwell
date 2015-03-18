package card
import (
    "orwell/orlib/sig"
    "time"
    "encoding/asn1"
    "errors"
)

type Card struct {
    Key sig.PubKey
    Version int64
    Expires time.Time
    Records []Record
    Signature []byte
}

type Asn1Card struct {
    Key []byte // The standard does not specify key structure, future expansion possible (move to PublicKeySpecInfo?)
    Payload   Asn1Payload
    Signature []byte
}

type Asn1Payload struct {
    Version int64
    Expires int64
    Records []Record
}

type Record struct {
    Key string `asn1:"utf8"`
    Type string `asn1:"utf8"`
    Value string `asn1:"utf8"`
}

func MarshalPayload(card *Card) ([]byte, error) {
    if card == nil {
        return nil, errors.New("Nil card value provided")
    }
    // TODO: inspect the consequences of int64/uint64 overflow
    payload := Asn1Payload{}
    payload.Version = card.Version
    payload.Expires = card.Expires.Unix()
    payload.Records = card.Records
    return asn1.Marshal(payload)
}

func Marshal(card *Card) ([]byte, error) {
    c := Asn1Card {}
    c.Key = card.Key.Serialize()
    c.Payload.Version = card.Version
    c.Payload.Expires = card.Expires.Unix()
    c.Payload.Records = card.Records
    c.Signature = card.Signature
    return asn1.Marshal(c)
}

func Unmarshal(data []byte) (*Card, error) {
    c := Asn1Card{}
    rest, err := asn1.Unmarshal(data, &c)
    if len(rest) > 0 {
        return nil, errors.New("Serialized card too long (bytes remaining)")
    }
    if err != nil {
        return nil, err
    }
    card := Card{}
    card.Key, err = sig.ParsePubKey(c.Key)
    if err != nil {
        return nil, err
    }
    card.Version = c.Payload.Version
    card.Expires = time.Unix(c.Payload.Expires, 0)
    card.Records = c.Payload.Records
    card.Signature = c.Signature
    if Verify(&card) {
        return &card, nil
    } else {
        return nil, errors.New("Card verification failed")
    }
}

func Sign(card *Card, key sig.PrvKey) error {
    payload, err := MarshalPayload(card)
    if err != nil {
        return err
    } else {
        card.Signature = key.Sign(payload)
        card.Key = key.PublicPart()
        return nil
    }
}

func Verify(card *Card) bool {
    payload, err := MarshalPayload(card)
    if err != nil {
        return false
    } else {
        return card.Key.Verify(payload, card.Signature)
    }
}
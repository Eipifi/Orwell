package sig

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
    "errors"
)

type libPrvKey ecdsa.PrivateKey
type libPubKey ecdsa.PublicKey

type asn1EcdsaPrvKey struct {
	X, Y, D *big.Int
}

type asn1EcdsaPubKey struct {
	X, Y *big.Int
}

type asn1EcdsaSignature struct {
	R, S *big.Int
}

func NewPrivateKey() PrvKey {
	var curve elliptic.Curve = elliptic.P256()
	key, _ := ecdsa.GenerateKey(curve, rand.Reader)
	return libPrvKey(*key)
}

func (k libPrvKey) PublicPart() PubKey {
	return libPubKey(k.PublicKey)
}

func (k libPrvKey) Sign(data []byte) []byte {
	r, s, _ := ecdsa.Sign(rand.Reader, (*ecdsa.PrivateKey)(&k), HashSlice(data))
	sequence := asn1EcdsaSignature{r, s}
	result, _ := asn1.Marshal(sequence)
	return result
}

func (k libPubKey) Verify(data []byte, signature []byte) bool {
	sequence := new(asn1EcdsaSignature)
	rest, err := asn1.Unmarshal(signature, sequence)

	if (len(rest) != 0) || (err != nil) {
		return false
	} else {
		return ecdsa.Verify((*ecdsa.PublicKey)(&k), HashSlice(data), sequence.R, sequence.S)
	}
}

func (k libPubKey) Marshal() []byte {
	sequence := asn1EcdsaPubKey{k.X, k.Y}
	result, _ := asn1.Marshal(sequence)
	return result
}

func (k libPrvKey) Marshal() []byte {
	sequence := asn1EcdsaPrvKey{k.PublicKey.X, k.PublicKey.Y, k.D}
	result, _ := asn1.Marshal(sequence)
	return result
}

func UnmarshalPubKey(data []byte) (PubKey, error) {
    var qwe asn1EcdsaPubKey
    rest, err := asn1.Unmarshal(data, &qwe)
    if len(rest) != 0 {
        return nil, errors.New("Bytes remaining after deserializing public key")
    }
    if err != nil {
        return nil, err
    }
    var res libPubKey
    res.Curve = elliptic.P256()
    res.X = qwe.X
    res.Y = qwe.Y
    return res, nil
}

func (k libPubKey) Id() ID {
	return Hash(k.Marshal())
}

func (k libPrvKey) Id() ID {
	return k.PublicPart().Id()
}

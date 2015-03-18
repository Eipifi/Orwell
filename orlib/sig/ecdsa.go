package sig
import (
    "crypto/ecdsa"
    "encoding/asn1"
    "crypto/rand"
    "crypto/elliptic"
    "crypto/x509"
    "errors"
    "math/big"
)

type asn1EcdsaSignature struct {
    R, S *big.Int
}

type lPubKey struct {
    key *ecdsa.PublicKey
}

type lPrvKey struct {
    key *ecdsa.PrivateKey
}

func NewPrvKey() PrvKey {
    var curve elliptic.Curve = elliptic.P256()
    key, _ := ecdsa.GenerateKey(curve, rand.Reader)
    return lPrvKey{key}
}

func (k lPubKey) Id() ID {
    return Hash(k.Marshal())
}

func (k lPrvKey) Id() ID {
    return k.PublicPart().Id()
}

func (k lPrvKey) PublicPart() PubKey {
    return lPubKey{&k.key.PublicKey}
}

func (k lPrvKey) Sign(data []byte) []byte {
    r, s, _ := ecdsa.Sign(rand.Reader, k.key, HashSlice(data))
    sequence := asn1EcdsaSignature{r, s}
    result, _ := asn1.Marshal(sequence)
    return result
}

func (k lPubKey) Verify(data []byte, signature []byte) bool {
    sequence := new(asn1EcdsaSignature)
    rest, err := asn1.Unmarshal(signature, sequence)

    if (len(rest) != 0) || (err != nil) {
        return false
    } else {
        return ecdsa.Verify(k.key, HashSlice(data), sequence.R, sequence.S)
    }
}

func (k lPubKey) Marshal() []byte {
    data, _ := x509.MarshalPKIXPublicKey(k.key)
    return data
}

func (k lPrvKey) Marshal() []byte {
    data, _ := x509.MarshalECPrivateKey(k.key)
    return data
}

func UnmarshalPK(data []byte) (PubKey, error) {
    pub, err := x509.ParsePKIXPublicKey(data)
    if err != nil {
        return nil, err
    } else {
        switch key := pub.(type) {
            case *ecdsa.PublicKey:
                return lPubKey{key}, nil
            default:
                return nil, errors.New("x509: only RSA and ECDSA public keys supported")
        }
    }
}
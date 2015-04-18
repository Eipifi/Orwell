package sig
import (
    "crypto/ecdsa"
    "crypto/x509"
    "encoding/asn1"
    "math/big"
    "crypto/rand"
    "crypto/elliptic"
)

// ECDSA Public Key

type ecdsaPubKey struct {
    obj *ecdsa.PublicKey
}

type ecdsaSignature struct {
    R, S *big.Int
}

func (key ecdsaPubKey) Serialize() []byte {
    data, _ := x509.MarshalPKIXPublicKey(key.obj)
    return data
}

func (key ecdsaPubKey) Id() *ID {
    return Hash(key.Serialize())
}

func (key ecdsaPubKey) Verify(data []byte, signature []byte) bool {
    sequence := ecdsaSignature{}
    rest, err := asn1.Unmarshal(signature, &sequence)
    if (len(rest) != 0) || (err != nil) {
        return false
    } else {
        return ecdsa.Verify(key.obj, HashSlice(data), sequence.R, sequence.S)
    }
}

// ECDSA Private Key

type ecdsaPrvKey struct {
    obj *ecdsa.PrivateKey
}

func (key ecdsaPrvKey) Serialize() []byte {
    data, _ := x509.MarshalECPrivateKey(key.obj)
    return data
}

func (key ecdsaPrvKey) PublicPart() PubKey {
    return ecdsaPubKey{&key.obj.PublicKey}
}

func (key ecdsaPrvKey) Sign(data []byte) []byte {
    r, s, _ := ecdsa.Sign(rand.Reader, key.obj, HashSlice(data))
    sequence := ecdsaSignature{r, s}
    result, _ := asn1.Marshal(sequence)
    return result
}

// Other

func NewEcdsaPrvKey() PrvKey {
    key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    return ecdsaPrvKey{key}
}
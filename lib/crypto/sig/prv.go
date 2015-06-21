package sig
import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "orwell/lib/crypto/hash"
    "crypto/x509"
)

var ENTROPY_SOURCE = rand.Reader
type PrvKey ecdsa.PrivateKey

func (k *PrvKey) ReadBytes(derBytes []byte) error {
    prv, err := x509.ParseECPrivateKey(derBytes)
    ptr := (*PrvKey)(prv)
    *k = *ptr
    return err
}

func (k *PrvKey) WriteBytes() ([]byte, error) {
    ptr := (*ecdsa.PrivateKey)(k)
    return x509.MarshalECPrivateKey(ptr)
}

func Create() (*PrvKey, error) {
    ptr, err := ecdsa.GenerateKey(elliptic.P256(), ENTROPY_SOURCE)
    return (*PrvKey)(ptr), err
}

func (k *PrvKey) PublicPart() *PubKey {
    ptr := &(k.PublicKey)
    return (*PubKey)(ptr)
}

func (k *PrvKey) Sign(payload []byte) (*Signature, error) {
    ptr := (*ecdsa.PrivateKey)(k)
    r, s, err := ecdsa.Sign(ENTROPY_SOURCE, ptr, hash.HashBytes(payload))
    if err != nil { return nil, err }
    return &Signature{r, s}, nil
}
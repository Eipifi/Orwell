package sig
import (
    "crypto/ecdsa"
    "crypto/x509"
    "crypto/rand"
    "crypto/elliptic"
    "orwell/orlib/butils"
    "orwell/orlib/crypto/hash"
)

type PrivateKey struct {
    obj *ecdsa.PrivateKey
}

func (k *PrivateKey) ReadBytes(data []byte) error {
    prv, err := x509.ParseECPrivateKey(data)
    k.obj = prv
    return err
}

func (k *PrivateKey) WriteBytes() ([]byte, error) {
    return x509.MarshalECPrivateKey(k.obj)
}

func (k *PrivateKey) PublicPart() *PublicKey {
    return &PublicKey{&k.obj.PublicKey}
}

func (k *PrivateKey) Sign(payload []byte) (*Signature, error) {
    h := hash.Hash(payload)
    r, s, err := ecdsa.Sign(rand.Reader, k.obj, h)
    if err != nil { return nil, err }
    return &Signature{r, s}, nil
}

func (k *PrivateKey) SignByteWritable(w butils.ByteWritable) (*Signature, error) {
    buf, err := w.WriteBytes()
    if err != nil { return nil, err }
    return k.Sign(buf)
}

func CreateKey() (*PrivateKey, error) {
    key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil { return nil, err }
    return &PrivateKey{key}, nil
}
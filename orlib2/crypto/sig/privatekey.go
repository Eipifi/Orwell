package sig
import (
    "crypto/ecdsa"
    "io"
    "crypto/x509"
    "crypto/rand"
    "orwell/orlib/protocol/types"
    "crypto/elliptic"
    "orwell/orlib2/butils"
)

// TODO: make PrivateKey implement Readable

type PrivateKey struct {
    obj *ecdsa.PrivateKey
}

func (k *PrivateKey) ReadBytes(data []byte) error {
    prv, err := x509.ParseECPrivateKey(data)
    if err != nil { return err }
    // TODO: validate if the proper elliptic curve is used
    k.obj = prv
    return
}

func (k *PrivateKey) Write(w io.Writer) error {
    data, err := x509.MarshalECPrivateKey(k.obj)
    if err != nil { return err }
    return w.Write(data)
}

func (k *PrivateKey) PublicPart() *PublicKey {
    return &PublicKey{&k.obj.PublicKey}
}

func (k *PrivateKey) Sign(payload []byte) (*Signature, error) {
    h := types.HashSlice(payload)
    r, s, err := ecdsa.Sign(rand.Reader, k.obj, h)
    if err != nil { return nil, err }
    return &Signature{r, s}, nil
}

func (k *PrivateKey) SignWritable(w butils.Writable) (*Signature, err error) {
    var buf []byte
    if buf, err = butils.WriteToBytes(w); err != nil { return }
    return k.Sign(buf)
}

func CreateKey() (*PrivateKey, error) {
    key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil { return }
    return &PrivateKey{key}
}
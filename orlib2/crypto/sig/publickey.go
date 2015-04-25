package sig
import (
    "crypto/ecdsa"
    "io"
    "crypto/x509"
    "errors"
    "orwell/orlib/protocol/types"
)

// TODO: make PublicKey implement Readable

type PublicKey struct {
    obj *ecdsa.PublicKey
}

func (k *PublicKey) Write(w io.Writer) error {
    data, err := x509.MarshalPKIXPublicKey(k.obj)
    if err != nil { return err }
    return w.Write(data)
}

func (k *PublicKey) ReadBytes(data []byte) error {
    pub, err := x509.ParsePKIXPublicKey(data)
    if err != nil { return err }
    switch pub := pub.(type) {
        case *ecdsa.PublicKey:
            k.obj = pub
        default:
            return errors.New("Unsupported key type")
    }
    return
}

func (k *PublicKey) WriteBytes() (data []byte, err error) {
    return x509.MarshalPKIXPublicKey(k.obj)
}

func (k *PublicKey) Verify(payload []byte, signature *Signature) error {
    h := types.HashSlice(payload)
    if ecdsa.Verify(k.obj, h, signature.R, signature.S) {
        return nil
    } else {
        return errors.New("Signature verification failed")
    }
}
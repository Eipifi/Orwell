package sig
import (
    "crypto/ecdsa"
    "crypto/x509"
    "errors"
    "orwell/lib/crypto/hash"
    "io"
    "orwell/lib/butils"
    "orwell/lib/foo"
    "orwell/lib/utils"
)

const PUBKEY_MAX_LENGTH uint64 = 4096

type PubKey ecdsa.PublicKey

func (k *PubKey) Read(r io.Reader) error {
    buf, err := butils.ReadVarBytes(r, PUBKEY_MAX_LENGTH)
    if err != nil { return err }
    val, err := x509.ParsePKIXPublicKey(buf)
    if err != nil { return err }
    if pub, ok := val.(*ecdsa.PublicKey); ok {
        ptr := (*PubKey)(pub)
        *k = *ptr
        return nil
    }
    return errors.New("Invalid public key type")
}

func (k *PubKey) Write(w io.Writer) error {
    buf, err := x509.MarshalPKIXPublicKey((*ecdsa.PublicKey)(k))
    if err != nil { return err }
    return butils.WriteVarBytes(w, buf, PUBKEY_MAX_LENGTH)
}

func (k *PubKey) Verify(payload []byte, s *Signature) error {
    ptr := (*ecdsa.PublicKey)(k)
    if ecdsa.Verify(ptr, hash.HashBytes(payload), s.R, s.S) { return nil }
    return errors.New("Invalid signature")
}

func (k *PubKey) TryID() (id foo.U256, err error) {
    buf, err := butils.WriteToBytes(k)
    if err == nil { id = hash.Hash(buf) }
    return
}

func (k *PubKey) ID() foo.U256 {
    id, err := k.TryID()
    utils.Ensure(err)
    return id
}
package orchain
import (
    "orwell/lib/crypto/sig"
    "io"
    "orwell/lib/butils"
)

type Proof struct {
    PublicKey sig.PubKey
    Signature sig.Signature
}

func (p *Proof) Read(r io.Reader) (err error) {
    if err = p.PublicKey.Read(r); err != nil { return }
    if err = p.Signature.Read(r); err != nil { return }
    return
}

func (p *Proof) Write(w io.Writer) (err error) {
    if err = p.PublicKey.Write(w); err != nil { return }
    if err = p.Signature.Write(w); err != nil { return }
    return
}

func (p *Proof) Check(data []byte) error {
    return p.PublicKey.Verify(data, &(p.Signature))
}

func (p *Proof) CheckObject(w butils.Writable) error {
    data, err := butils.WriteToBytes(w)
    if err != nil { return err }
    return p.Check(data)
}
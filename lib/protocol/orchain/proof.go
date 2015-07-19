package orchain
import (
    "orwell/lib/crypto/sig"
    "io"
    "orwell/lib/foo"
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

func (p *Proof) Check(hash foo.U256) error {
    return p.PublicKey.Verify(hash[:], &(p.Signature))
}
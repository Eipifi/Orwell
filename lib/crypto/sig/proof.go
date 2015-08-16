package sig
import (
    "io"
    "orwell/lib/butils"
    "bytes"
)

type Proof struct {
    PublicKey PubKey
    Signature Signature
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

func (p *Proof) CheckWritable(w butils.Writable) error {
    data, err := butils.WriteToBytes(w)
    if err != nil { return err }
    return p.Check(data)
}

func (p *Proof) Sign(data []byte, key *PrvKey) (err error) {
    pk := key.PublicPart()
    sg, err := key.Sign(data)
    if err != nil { return }
    p.Signature = *sg
    p.PublicKey = *pk
    return
}

func (p *Proof) SignWritable(w butils.Writable, key *PrvKey) (err error) {
    data, err := butils.WriteToBytes(w)
    if err != nil { return }
    return p.Sign(data, key)
}

/////////////////////////////////////////////////////////////

type HeadWritable interface {
    WriteHead(w io.Writer) error
}

func (p *Proof) SignHead(w HeadWritable, key *PrvKey) (err error) {
    buf := &bytes.Buffer{}
    if err = w.WriteHead(buf); err != nil { return }
    return p.Sign(buf.Bytes(), key)
}

func (p *Proof) CheckHead(w HeadWritable) (err error) {
    buf := &bytes.Buffer{}
    if err = w.WriteHead(buf); err != nil { return }
    return p.Check(buf.Bytes())
}

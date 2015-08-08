package orchain
import "io"

type Transfer struct {
    Domain Domain
    Proof Proof
}

func (t *Transfer) Read(r io.Reader) (err error) {
    if err = t.Domain.Read(r); err != nil { return }
    if err = t.Proof.Read(r); err != nil { return }
    return
}

func (t *Transfer) Write(w io.Writer) (err error) {
    if err = t.Domain.Write(w); err != nil { return }
    if err = t.Proof.Write(w); err != nil { return }
    return
}
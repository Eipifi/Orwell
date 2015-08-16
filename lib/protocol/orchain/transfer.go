package orchain
import (
    "io"
    "orwell/lib/foo"
    "errors"
    "fmt"
    "orwell/lib/crypto/sig"
)

type Transfer struct {
    Domain Domain
    Proof sig.Proof
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

func (t *Transfer) Verify(owner foo.U256) (err error) {
    if t.Domain.Owner != owner { return errors.New("Invalid domain owner") }
    return t.Proof.CheckWritable(&t.Domain)
}

func (t *Transfer) String() string {
    return fmt.Sprintf("Transfer %v", t.Domain)
}
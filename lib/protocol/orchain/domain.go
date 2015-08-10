package orchain
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
)

const DOMAIN_MAX_LENGTH uint64 = 32

type Domain struct {
    Name string // TODO: perform lexical checks, ensure no special chars and whitespaces
    Owner foo.U256
    ValidUntil uint64
}

func (d *Domain) Read(r io.Reader) (err error) {
    if d.Name, err = butils.ReadString(r, DOMAIN_MAX_LENGTH); err != nil { return }
    if err = d.Owner.Read(r); err != nil { return }
    if d.ValidUntil, err = butils.ReadUint64(r); err != nil { return }
    return
}

func (d *Domain) Write(w io.Writer) (err error) {
    if err = butils.WriteString(w, d.Name, DOMAIN_MAX_LENGTH); err != nil { return }
    if err = d.Owner.Write(w); err != nil { return }
    if err = butils.WriteUint64(w, d.ValidUntil); err != nil { return }
    return
}
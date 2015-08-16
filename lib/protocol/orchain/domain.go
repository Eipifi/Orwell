package orchain
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
    "orwell/lib/utils"
)

const DOMAIN_MAX_LENGTH uint64 = 32

type Domain struct {
    Name string
    Owner foo.U256
    ValidUntilBlock uint64
}

func (d *Domain) Read(r io.Reader) (err error) {
    if err = CheckDomainString(d.Name); err != nil { return }
    if d.Name, err = butils.ReadString(r, DOMAIN_MAX_LENGTH); err != nil { return }
    if err = d.Owner.Read(r); err != nil { return }
    if d.ValidUntilBlock, err = butils.ReadUint64(r); err != nil { return }
    return
}

func (d *Domain) Write(w io.Writer) (err error) {
    if err = CheckDomainString(d.Name); err != nil { return }
    if err = butils.WriteString(w, d.Name, DOMAIN_MAX_LENGTH); err != nil { return }
    if err = d.Owner.Write(w); err != nil { return }
    if err = butils.WriteUint64(w, d.ValidUntilBlock); err != nil { return }
    return
}

func (d *Domain) TryID() (foo.U256, error) {
    return hash.HashOf(d)
}

func (d *Domain) ID() foo.U256 {
    id, err := d.TryID()
    utils.Ensure(err)
    return id
}

func CheckDomainString(domain string) error {
    // TODO: perform lexical checks, ensure no special chars and whitespaces
    return nil
}
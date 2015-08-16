package orcache
import (
    "orwell/lib/crypto/sig"
    "io"
    "orwell/lib/butils"
    "orwell/lib/foo"
)

const MAX_ENTRIES uint64 = 128

type Card struct {
    Version uint64
    Timestamp uint64
    Entries []Entry
    Proof sig.Proof
}

func (c *Card) ReadHead(r io.Reader) (err error) {
    if c.Version, err = butils.ReadVarUint(r); err != nil { return }
    if c.Timestamp, err = butils.ReadUint64(r); err != nil { return }
    if err = butils.ReadSlice(r, MAX_ENTRIES, &c.Entries); err != nil { return }
    return
}

func (c *Card) WriteHead(w io.Writer) (err error) {
    if err = butils.WriteVarUint(w, c.Version); err != nil { return }
    if err = butils.WriteUint64(w, c.Timestamp); err != nil { return }
    if err = butils.WriteSlice(w, MAX_ENTRIES, c.Entries); err != nil { return }
    return
}

func (c *Card) Read(r io.Reader) (err error) {
    if err = c.ReadHead(r); err != nil { return }
    if err = c.Proof.Read(r); err != nil { return }
    return
}

func (c *Card) Write(w io.Writer) (err error) {
    if err = c.WriteHead(w); err != nil { return }
    if err = c.Proof.Write(w); err != nil { return }
    return
}

func (c *Card) ID() foo.U256 {
    return c.Proof.PublicKey.ID()
}


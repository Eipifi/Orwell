package common
import (
    "io"
    "orwell/orlib/butils"
)

type Address struct {
    IP [16]byte
    Port uint16
    Nonce uint64
}

func (a *Address) Read(r io.Reader) (err error) {
    if err = butils.ReadFull(r, a.IP[:]); err != nil { return }
    if a.Port, err = butils.ReadUint16(r); err != nil { return }
    a.Nonce, err = butils.ReadUint64(r)
    return
}

func (a *Address) Write(w io.Writer) (err error) {
    if err = butils.WriteFull(w, a.IP[:]); err != nil { return }
    if err = butils.WriteUint16(w, a.Port); err != nil { return }
    return butils.WriteUint64(w, a.Nonce)
}
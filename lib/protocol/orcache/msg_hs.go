package orcache
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
)

const MAX_ADDR_LENGTH uint64 = 128

type Handshake struct {
    ID foo.U256
    Address string
}

func (m *Handshake) Read(r io.Reader) (err error) {
    if err = m.ID.Read(r); err != nil { return }
    if m.Address, err = butils.ReadString(r, MAX_ADDR_LENGTH); err != nil { return }
    return
}

func (m *Handshake) Write(w io.Writer) (err error) {
    if err = m.ID.Write(w); err != nil { return }
    if err = butils.WriteString(w, m.Address, MAX_ADDR_LENGTH); err != nil { return }
    return
}

type HandshakeAck struct {}
func (*HandshakeAck) Read(io.Reader) error { return nil }
func (*HandshakeAck) Write(io.Writer) error { return nil }
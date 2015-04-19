package orcache
import (
    "orwell/orlib/comm"
    "orwell/orlib/protocol/types"
)

type Handshake struct {
    Magic uint32
    Version uint64
    UserAgent string
    Address *types.Address
}

func (m *Handshake) Read(r *comm.Reader) (err error) {
    if m.Magic, err = r.ReadUint32(); err != nil { return }
    if m.Version, err = r.ReadVaruint(); err != nil { return }
    if m.UserAgent, err = r.ReadStr(); err != nil { return }

    var f uint8
    if f, err = r.ReadUint8(); err != nil { return }
    if f & 0x01 > 0 {
        m.Address = &types.Address{}
        if m.Address.Read(r) != nil { return }
    }
    return
}

func (m *Handshake) Write(w *comm.Writer) {
    w.WriteUint32(m.Magic)
    w.WriteVaruint(m.Version)
    w.WriteString(m.UserAgent)
    if m.Address == nil {
        w.WriteUint8(0)
    } else {
        w.WriteUint8(1)
        m.Address.Write(w)
    }
}
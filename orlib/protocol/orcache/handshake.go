package orcache
import (
    "io"
    "orwell/orlib/butils"
    "orwell/orlib/protocol/common"
)


type Handshake struct {
    Magic uint32
    Version uint64
    UserAgent string
    Address *common.Address
}

func (m *Handshake) Read(r io.Reader) (err error) {
    if m.Magic, err = butils.ReadUint32(r); err != nil { return }
    if m.Version, err = butils.ReadVarUint(r); err != nil { return }
    if m.UserAgent, err = butils.ReadString(r); err != nil { return }

    var f uint8
    if f, err = butils.ReadUint8(r); err != nil { return }
    if f & 0x01 > 0 {
        m.Address = &common.Address{}
        if m.Address.Read(r) != nil { return }
    }
    return
}

func (m *Handshake) Write(w io.Writer) (err error) {
    if err = butils.WriteUint32(w, m.Magic); err != nil { return }
    if err = butils.WriteVarUint(w, m.Version); err != nil { return }
    if err = butils.WriteString(w, m.UserAgent); err != nil { return }
    if m.Address == nil {
        return butils.WriteByte(w, 0)
    } else {
        if err = butils.WriteByte(w, 1); err != nil { return }
        return m.Address.Write(w)
    }
}
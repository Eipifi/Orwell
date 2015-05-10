package orcache
import (
    "io"
    "orwell/orlib/butils"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
)


type Handshake struct {
    Magic uint32
    Version uint64
    UserAgent string
    Port common.Port
    ID *hash.ID
}

func (*Handshake) Code() byte { return 0x01 }

func (m *Handshake) Read(r io.Reader) (err error) {
    if m.Magic, err = butils.ReadUint32(r); err != nil { return }
    if m.Version, err = butils.ReadVarUint(r); err != nil { return }
    if m.UserAgent, err = butils.ReadString(r); err != nil { return }
    var port uint16
    if port, err = butils.ReadUint16(r); err != nil { return }
    m.Port = common.Port(port)

    var flag uint8
    if flag, err = butils.ReadUint8(r); err != nil { return }
    if flag & 0x01 > 0 {
        m.ID = &hash.ID{}
        if err = m.ID.Read(r); err != nil { return }
    }
    return
}

func (m *Handshake) Write(w io.Writer) (err error) {
    if err = butils.WriteUint32(w, m.Magic); err != nil { return }
    if err = butils.WriteVarUint(w, m.Version); err != nil { return }
    if err = butils.WriteString(w, m.UserAgent); err != nil { return }
    if err = butils.WriteUint16(w, uint16(m.Port)); err != nil { return }
    if m.ID == nil {
        return butils.WriteUint8(w, 0x00)
    } else {
        if err = butils.WriteUint8(w, 0x01); err != nil { return }
        return m.ID.Write(w)
    }
}
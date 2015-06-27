package orchain
import (
    "io"
    "orwell/lib/butils"
)

// TODO: map size limits

const MSG_CODE_HANDSHAKE_REQ = 0x00
const MSG_CODE_HANDSHAKE_RSP = 0x80

type HandshakeReq struct {
    Magic uint32
    Fields map[string] string
}

func (m *HandshakeReq) Read(r io.Reader) (err error) {
    if m.Magic, err = butils.ReadUint32(r); err != nil { return }
    var num uint64
    if num, err = butils.ReadVarUint(r); err != nil { return }
    m.Fields = make(map[string] string)
    for i := 0; i < int(num); i += 1 {
        var k, v string
        if k, err = butils.ReadString(r); err != nil { return }
        if v, err = butils.ReadString(r); err != nil { return }
        m.Fields[k] = v
    }
    return
}

func (m *HandshakeReq) Write(w io.Writer) (err error) {
    if err = butils.WriteUint32(w, m.Magic); err != nil { return }
    if err = butils.WriteVarUint(w, uint64(len(m.Fields))); err != nil { return }
    for k, v := range m.Fields {
        if err = butils.WriteString(w, k); err != nil { return }
        if err = butils.WriteString(w, v); err != nil { return }
    }
    return nil
}

type HandshakeRsp struct {}

func (m *HandshakeRsp) Read(io.Reader) error { return nil }
func (m *HandshakeRsp) Write(io.Writer) error { return nil }

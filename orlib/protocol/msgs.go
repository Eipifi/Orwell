package protocol
import (
    "errors"
)

/*
HELLO		: magic, pv, [id], [address]
HELLO_ACK	: magic, pv, id
GET			: token, ttl, id, version
CARD		: token, card
NOPE		: token, ttl
PUT			: token, ttl, card
SAVED		: token, ttl
NEIGHBOURS	: id

*/

type Message interface {
    Command() uint64
    WriteTo(w *Writer)
}

type Frame struct {
    Command uint64
    Payload []byte
}

// <Handshake> -------------------------------------------------------------------------
const CodeHandshake     = 0x01

type Handshake struct {
    Magic uint32
    Version uint64
    UserAgent string
    Address *Address
}

func (m *Handshake) Command() { return CodeHandshake }

func (m *Handshake) WriteTo(w *Writer) {
    w.WriteUint32(m.Magic)
    w.WriteVaruint(m.Version)
    w.WriteString(m.UserAgent)
    if m.Address == nil {
        w.WriteUint8(0)
    } else {
        w.WriteUint8(1)
        w.WriteAddress(m.Address)
    }
}

func (r *Reader) ReadHandshake() (msg *Handshake, err error) {
    msg = &Handshake{}
    if msg.Magic, err = r.ReadUint32(); err != nil { return }
    if msg.Version, err = r.ReadVaruint(); err != nil { return }
    if msg.UserAgent, err = r.ReadStr(); err != nil { return }

    var f uint8
    if f, err = r.ReadUint8(); err != nil { return }
    if f & 0x01 > 0 {
        if msg.Address, err = r.ReadAddress(); err != nil { return }
    }
    return
}

// </Handshake> ------------------------------------------------------------------------

func (f *Frame) Unpack() (interface{}, error) {
    r := NewBytesReader(f.Payload)
    switch f.Command {
        case CodeHandshake:
            return r.ReadHandshake()
        default:
            return nil, errors.New("Unknown command code")
    }
}
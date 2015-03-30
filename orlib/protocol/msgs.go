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
const CodeHandshake = 0x01

type Handshake struct {
    Magic uint32
    Version uint64
    UserAgent string
    Address *Address
}

func (m *Handshake) Command() uint64 { return CodeHandshake }

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

func (r *Reader) ReadHandshake(msg *Handshake) (err error) {
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
// <HandshakeAck> ----------------------------------------------------------------------
const CodeHandshakeAck = 0x81

type HandshakeAck struct { }

func (m *HandshakeAck) Command() uint64 { return CodeHandshakeAck }

func (m *HandshakeAck) WriteTo(w * Writer) { }

func (r *Reader) ReadHandshakeAck(*HandshakeAck) (error) { return nil }

// </HandshakeAck> ---------------------------------------------------------------------

func (r *Reader) ReadAnyMessage() (m Message, err error) {
    var f *Frame
    if f, err = r.ReadFrame(); err != nil { return }
    return Switch1(NewBytesReader(f.Payload), f.Command)
}

func (r *Reader) ReadMessage(m Message) (err error) {
    var f *Frame
    if f, err = r.ReadFrame(); err != nil { return }
    if f.Command != m.Command() { errors.New("Unexpected frame command") }
    return Switch2(NewBytesReader(f.Payload), m)
}

// TODO: unify these two switch statements

func Switch1(r *Reader, command uint64) (m Message, err error) {
    switch command {
        case CodeHandshake:
            msg := Handshake{}
            err = r.ReadHandshake(&msg)
            return &msg, err
        case CodeHandshakeAck:
            msg := HandshakeAck{}
            err = r.ReadHandshakeAck(&msg)
            return &msg, err
        default:
            return nil, errors.New("Unknown frame command")
    }
}

func Switch2(r *Reader, m Message) error {
    switch m.Command() {
        case CodeHandshake:
            return r.ReadHandshake(m.(*Handshake))
        case CodeHandshakeAck:
            return r.ReadHandshakeAck(m.(*HandshakeAck))
        default:
            return errors.New("Unknown frame command")
    }
}
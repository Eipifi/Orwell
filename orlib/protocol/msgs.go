package protocol
import (
    "reflect"
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

type Msg interface {
    Read(r *Reader) error
    Write(w *Writer)
}

type Frame struct {
    Command uint64
    Payload []byte
}

type msgTypeEntry struct {
    Command uint64
    Type reflect.Type
}

var msgTypes = []msgTypeEntry {
    msgTypeEntry{0x01, reflect.TypeOf(Handshake{})},
}

func GetMsgCommand(m Msg) uint64 {
    t := reflect.TypeOf(m)
    for _, e := range msgTypes {
        if e.Type == t {
            return e.Command
        }
    }
    return 0x00
}

func GetCommandMsg(c uint64) Msg {
    for _, e := range msgTypes {
        if e.Command == c {
            v := reflect.New(e.Type)
            return v.Interface().(Msg)
        }
    }
    return nil
}

func (w *Writer) WriteFramedMessage(m Msg) {
    v := NewWriter()
    m.Write(v)
    w.WriteVaruint(GetMsgCommand(m))
    w.Write(v.Peek())
}

func (r *Reader) ReadFramedMessage() (m Msg, err error) {
    var f *Frame
    if f, err = r.ReadFrame(); err != nil { return }
    if m = GetCommandMsg(f.Command); m == nil { return nil, errors.New("Unrecognized message type") }
    if err = m.Read(NewBytesReader(f.Payload)); err != nil { return }
    return
}

func (r *Reader) ReadSpecificFramedMessage(m Msg) (err error) {
    var f *Frame
    if f, err = r.ReadFrame(); err != nil { return }
    if f.Command != GetMsgCommand(m) { return errors.New("Unexpected message type") }
    return m.Read(NewBytesReader(f.Payload))
}

///////////////////////////////////////////////////////////////////////////

type Handshake struct {
    Magic uint32
    Version uint64
    UserAgent string
    Address *Address
}

func (m *Handshake) Read(r *Reader) (err error) {
    if m.Magic, err = r.ReadUint32(); err != nil { return }
    if m.Version, err = r.ReadVaruint(); err != nil { return }
    if m.UserAgent, err = r.ReadStr(); err != nil { return }

    var f uint8
    if f, err = r.ReadUint8(); err != nil { return }
    if f & 0x01 > 0 {
        if m.Address, err = r.ReadAddress(); err != nil { return }
    }
    return
}

func (m *Handshake) Write(w *Writer) {
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

///////////////////////////////////////////////////////////////////////////

type HandshakeAck struct { }

func (m *HandshakeAck) Read(r *Reader) error { return nil }

func (m *HandshakeAck) Write(w *Writer) { }

///////////////////////////////////////////////////////////////////////////

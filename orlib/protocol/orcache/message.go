package orcache
import (
    "orwell/orlib"
    "io"
    "orwell/orlib/butils"
    "errors"
)

type Message struct {
    Command uint64
    Chunk butils.Chunk
}

// TODO: maybe remove the length field (it is redundant, after all)

func (m *Message) Read(r io.Reader) (err error) {
    if m.Command, err = butils.ReadVarUint(r); err != nil { return }
    m.Chunk = commandToChunk(m.Command)
    if m.Chunk == nil { return errors.New("Unrecognized message code") }
    var payload []byte
    if payload, err = butils.ReadVarBytes(r); err != nil { return }
    return butils.ReadAllInto(m.Chunk, payload)
}

func (m *Message) Write(w io.Writer) (err error) {
    if err = butils.WriteVarUint(w, m.Command); err != nil { return }
    var payload []byte
    if payload, err = butils.WriteToBytes(m.Chunk); err != nil { return }
    return butils.WriteVarBytes(w, payload)
}

func NewMessage(chunk butils.Chunk) *Message {
    m := &Message{}
    m.Chunk = chunk
    m.Command = chunkToCommand(chunk)
    if m.Command == 0 { return nil }
    return m
}

func commandToChunk(command uint64) butils.Chunk {
    if command == 0x01 { return &Handshake{} }
    if command == 0x02 { return &Get{} }
    if command == 0x03 { return &Publish{} }
    if command == 0x81 { return &HandshakeAck{} }
    if command == 0x82 { return &CardFound{} }
    if command == 0x83 { return &CardNotFound{} }
    return nil
}

func chunkToCommand(chunk butils.Chunk) uint64 {
    switch chunk := chunk.(type) {
        case *Handshake:    return 0x01
        case *Get:          return 0x02
        case *Publish:      return 0x03
        case *HandshakeAck: return 0x81
        case *CardFound:    return 0x82
        case *CardNotFound: return 0x83
    }
    return 0
}
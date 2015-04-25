package orcache
import (
    "orwell/orlib2"
    "io"
    "orwell/orlib2/butils"
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
    return nil
}

func chunkToCommand(chunk butils.Chunk) uint64 {
    return 0
}
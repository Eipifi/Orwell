package orcache
import (
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
    var cmd uint64
    if cmd, err = butils.ReadVarUint(r); err != nil { return }

    if m.Chunk == nil { // if we do not expect any specific message type
        m.Command = cmd
        m.Chunk = commandToChunk(m.Command)
        if m.Chunk == nil { return errors.New("Unrecognized message code") }
    } else { // if we expect a specific type, specified by m.Chunk
        if m.Command != cmd { return errors.New("Unexpected message code") }
    }

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

func Msg(chunk butils.Chunk) *Message {
    m := &Message{}
    m.Chunk = chunk
    m.Command = chunkToCommand(chunk)
    if m.Command == 0 { return nil }
    return m
}

func commandToChunk(command uint64) butils.Chunk {
    if command == 0x01 { return &Handshake{} }
    if command == 0x81 { return &HandshakeAck{} }

    if command == 0x02 { return &FetchReq{} }
    if command == 0x82 { return &FetchRsp{} }

    if command == 0x03 { return &PublishReq{} }
    if command == 0x83 { return &PublishRsp{} }

    if command == 0x04 { return &PeersReq{} }
    if command == 0x84 { return &PeersRsp{} }

    return nil
}

func chunkToCommand(chunk butils.Chunk) uint64 {
    switch chunk.(type) {
        case *Handshake:        return 0x01
        case *HandshakeAck:     return 0x81

        case *FetchReq:         return 0x02
        case *FetchRsp:         return 0x82

        case *PublishReq:       return 0x03
        case *PublishRsp:       return 0x83

        case *PeersReq:         return 0x04
        case *PeersRsp:         return 0x84
    }
    return 0
}
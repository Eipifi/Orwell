package orcache
import (
    "net"
    "orwell/orlib/comm"
    "orwell/orlib/protocol/types"
    "errors"
)

type OrcacheMessenger struct {
    comm.Messenger
    hs *Handshake
}

func NewOrcacheMessenger(conn net.Conn, userAgent string, addr *types.Address) (ms *OrcacheMessenger, err error) {
    ms = &OrcacheMessenger{}
    ms.Messenger = *comm.NewMessager(conn, readSpecificFramedMessage, writeFramedMessage, readFramedMessage)

    // Send our Handshake
    if err = ms.Write(&Handshake{OrcacheMagic, SupportedVersion, userAgent, addr}); err != nil { return }

    // Await for the Handshake
    ms.hs = &Handshake{}
    if err = ms.Read(ms.hs); err != nil { return }

    // Send and await for ack
    if err = ms.Write(&HandshakeAck{}); err != nil { return }
    if err = ms.Read(&HandshakeAck{}); err != nil { return }
    return
}

func readSpecificFramedMessage(r *comm.Reader, m comm.Msg) (err error) {
    f := &types.Frame{}
    if err = f.Read(r); err != nil { return }
    if getMsgCommand(m) != f.Command { return errors.New("Unexpected message code") }
    return m.Read(comm.NewBytesReader(f.Payload))
}

func readFramedMessage(r *comm.Reader) (m comm.Msg, err error) {
    fr := &types.Frame{}
    if err = fr.Read(r); err != nil { return }
    if m = getCommandMsg(fr.Command); m == nil { return nil, errors.New("Unexpected message code") }
    return m, m.Read(comm.NewBytesReader(fr.Payload))
}

func writeFramedMessage(w *comm.Writer, m comm.Msg) {
    w2 := comm.NewWriter()
    m.Write(w2)
    f := &types.Frame{getMsgCommand(m), w2.Peek()}
    f.Write(w)
}
package orcache
import (
    "net"
    "orwell/orlib/comm"
    "orwell/orlib/protocol/types"
    "errors"
    "orwell/orlib/card"
)

type OrcacheMessenger struct {
    comm.Messenger
    hs *Handshake
}

func NewOrcacheMessenger(conn net.Conn, userAgent string, addr *types.Address) (ms *OrcacheMessenger, err error) {
    ms = &OrcacheMessenger{}
    ms.Messenger = *comm.NewMessenger(conn, readSpecificFramedMessage, writeFramedMessage, readFramedMessage)

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

func (ms *OrcacheMessenger) Get(id *types.ID, version uint64) (c *card.Card, err error) {
    t := types.RandomToken()
    if err = ms.Write(&Get{t, types.MaxTTLValue, id, version}); err != nil { return }
    var m comm.Msg
    if m, err = ms.ReadAny(); err != nil { return }
    switch m := m.(type) {
        case *CardFound:
            if m.Token == t {
                return card.Unmarshal(m.Card)
            } else {
                return nil, errors.New("Token mismatch")
            }
        case *CardNotFound:
        if m.Token == t {
            return nil, nil
        } else {
            return nil, errors.New("Token mismatch")
        }
        default:
            return nil, errors.New("Invalid message type received")
    }
}

func (ms *OrcacheMessenger) Put(c *card.Card) (ttl types.TTL, err error) {
    t := types.RandomToken()
    var raw []byte
    if raw, err = c.Marshal(); err != nil { return }
    if err = ms.Write(&Publish{t, types.MaxTTLValue, raw}); err != nil { return }
    p := &Published{}
    if err = ms.Read(p); err != nil { return }
    return p.TTL, nil
}

func SimpleClient(connectionString string) (ms *OrcacheMessenger, err error) {
    var c net.Conn
    if c, err = net.Dial("tcp", connectionString); err != nil { return }
    if ms, err = NewOrcacheMessenger(c, "", nil); err != nil { return }
    return
}
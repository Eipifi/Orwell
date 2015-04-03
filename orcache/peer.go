package main
import (
    "orwell/orlib/protocol"
    "net"
    "fmt"
    "errors"
)

const OrcacheMagic = 0x04cac11e
const SupportedVersion = 1
const UserAgent = "Orcache"

var connectionCounter uint64 = 0

type Peer struct {
    conn net.Conn
    hs *protocol.Handshake
    num uint64
    env *Env
}

func RunPeer(conn net.Conn, env *Env) {
    p := Peer{conn: conn, num: connectionCounter, env: env}
    connectionCounter += 1
    go func(){
        p.log(p.handle())
        p.log(p.conn.Close())
    }()
}

func (p *Peer) RelayGet(m *protocol.Get) <-chan *GetResponse {
    ch := make(chan *GetResponse, 1)
    /*
        if peer still active {
            put the get message on queue
        } else {
            ch <- nil
            close(ch)
        }
    */
    return ch
}

func (p *Peer) exchangeHandshakes() (err error) {
    // Initialize
    r := protocol.NewReader(p.conn)
    w := protocol.NewWriter()

    // Send our Handshake
    w.WriteFramedMessage(&protocol.Handshake{OrcacheMagic, SupportedVersion, UserAgent, nil})
    if err = w.Commit(p.conn); err != nil { return }

    // Await for the Handshake
    p.hs = &protocol.Handshake{}
    if err = r.ReadSpecificFramedMessage(p.hs); err != nil { return }

    // Send the HandshakeAck
    w.WriteFramedMessage(&protocol.HandshakeAck{})
    if err = w.Commit(p.conn); err != nil { return }

    // Await for the HandshakeAck
    var ack protocol.HandshakeAck
    if err = r.ReadSpecificFramedMessage(&ack); err != nil { return }
    return
}

func (p *Peer) handle() (err error) {
    if err = p.exchangeHandshakes(); err != nil { return }
    // TODO: handle messages
    return
}

func (p *Peer) log(err error) {
    if err != nil {
        fmt.Printf("Peer %d: %x\n", p.num, err)
    }
}
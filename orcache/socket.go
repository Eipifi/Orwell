package main
import (
    "net"
    "orwell/orlib/protocol"
)

const OrcacheMagic = 0x04cac11e
const SupportedVersion = 1
const UserAgent = "Orcache"

func exchangeHandshakes(conn net.Conn) (err error) {
    // Initialize
    r := protocol.NewReader(conn)
    w := protocol.NewWriter()

    // Send our Handshake
    w.WriteFramedMessage(&protocol.Handshake{OrcacheMagic, SupportedVersion, UserAgent, nil})
    if err = w.Commit(conn); err != nil { return }

    // Await for the Handshake
    var hs protocol.Handshake
    if err = r.ReadSpecificFramedMessage(&hs); err != nil { return }

    // Send the HandshakeAck
    w.WriteFramedMessage(&protocol.HandshakeAck{})
    if err = w.Commit(conn); err != nil { return }

    // Await for the HandshakeAck
    var ack protocol.HandshakeAck
    if err = r.ReadSpecificFramedMessage(&ack); err != nil { return }

    // Check the remote ID
    if hs.Address != nil {
        // TODO: inform connection manager about hs.Address.Id()
    }

    return // TODO: also return remote info
}

func socketHandler(conn net.Conn) (err error) {
    // Exchange handshakes
    if err = exchangeHandshakes(conn); err != nil { return }

    return
}


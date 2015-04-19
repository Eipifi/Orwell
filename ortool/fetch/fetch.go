package fetch
import (
    "net"
    "fmt"
    "orwell/orlib/protocol"
    "orwell/orlib/sig"
    "flag"
    "orwell/orlib/card"
)

const Usage = `usage: ortool fetch [--from ip:port] <address>

Asks the server for a card.

`
func Main(args []string) {
    if len(args) == 0 {
        fmt.Println("ID not provided.")
        return
    }

    var id *sig.ID
    var err error
    if id, err = sig.HexToID(args[len(args)-1]); err != nil {
        fmt.Println("Error: ", err)
        return
    }

    args = args[:len(args)-1]

    fs := flag.NewFlagSet("fetch", flag.ExitOnError)
    fSr := fs.String("from", "127.0.0.1:1984", "Orcache server address")
    fs.Parse(args)

    var conn net.Conn
    if conn, err = net.Dial("tcp", *fSr); err != nil {
        fmt.Println("Error:", err)
        return
    }

    r := protocol.NewReader(conn)
    w := protocol.NewWriter()

    w.WriteFramedMessage(&protocol.Handshake{protocol.OrcacheMagic, protocol.SupportedVersion, "Ortool", nil})
    w.WriteFramedMessage(&protocol.HandshakeAck{})
    w.WriteFramedMessage(&protocol.Get{protocol.RandomToken(), protocol.MaxTTLValue, id, 5})
    if err = w.Commit(conn); err != nil { return }

    if err = r.ReadSpecificFramedMessage(&protocol.Handshake{}); err != nil {
        fmt.Println("Failed to parse the handshake:", err)
        return
    }

    if err = r.ReadSpecificFramedMessage(&protocol.HandshakeAck{}); err != nil {
        fmt.Println("Failed to parse the handshake ack:", err)
        return
    }

    var msg protocol.Msg
    msg, err = r.ReadFramedMessage()
    if err != nil {
        fmt.Println("Failed to read response:", err)
        return
    }

    switch m := msg.(type) {
        case *protocol.CardFound:
            if c, e := card.Unmarshal(m.Card); e != nil {
                s := string(c.Payload.MarshalJSON())
                fmt.Println(s)
            } else {
                fmt.Println("Invalid card format or signature received.")
                return
            }
        case *protocol.CardNotFound:
            fmt.Printf("Failed to fetch the card. TTL=%d\n", m.TTL)
    }
}
package fetch
import (
    "net"
    "fmt"
    "orwell/orlib/protocol"
    "orwell/orlib/logging"
    "orwell/orlib/sig"
    "time"
)

const Usage = `usage: ortool fetch [--from ip:port] <address>

Asks the server for a card.

`
func Main(args []string) {
    err := main()
    if err != nil {
        fmt.Println(err)
    }
}

func main() (err error) {

    var conn net.Conn
    if conn, err = net.Dial("tcp", "127.0.0.1:1984"); err != nil { return }

    //r := protocol.NewReader(conn)
    w := protocol.NewWriter()

    w.WriteFramedMessage(&protocol.Handshake{0xcafebabe, 1, "", nil})
    w.WriteFramedMessage(&protocol.HandshakeAck{})
    w.WriteFramedMessage(&protocol.Get{65535, 16, sig.Hash([]byte("")), 5})
    if err = w.Commit(&logging.PrintWriter{conn}); err != nil { return }

    time.Sleep(time.Second * 5)

    return
}
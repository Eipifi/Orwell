package conv
import (
    "io"
    "orwell/orlib/protocol/common"
    "orwell/orlib/protocol/orcache"
)

const Magic uint32 = 0xcafebabe
const Version uint64 = 1

func ShakeHands(conn io.ReadWriter, userAgent string, address *common.Address) (hs *orcache.Handshake, err error) {
    hs = &orcache.Handshake{}
    if err = orcache.WriteMessage(conn, &orcache.Handshake{Magic, Version, userAgent, address}); err != nil { return }
    if err = orcache.ReadMessage(conn, hs); err != nil { return }
    if err = orcache.WriteMessage(conn, &orcache.HandshakeAck{}); err != nil { return }
    if err = orcache.ReadMessage(conn, &orcache.HandshakeAck{}); err != nil { return }
    return
}

func MessageListener(conn io.Reader) <-chan orcache.Message {
    c := make(chan orcache.Message)
    go func(){
        defer close(c)
        for {
            msg, err := orcache.ReadAnyMessage(conn)
            if err != nil { return }
            c <- msg
        }
    }()
    return c
}

func MessageSender(conn io.WriteCloser) chan<- orcache.Message {
    c := make(chan orcache.Message)
    go func(){
        defer conn.Close()
        for {
            msg, ok := <- c
            if ! ok { return }
            if orcache.WriteMessage(conn, msg) != nil { return }
        }
    }()
    return c
}
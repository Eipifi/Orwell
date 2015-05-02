package conv
import (
    "io"
    "orwell/orlib/protocol/common"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/butils"
)

const Magic uint32 = 0xcafebabe
const Version uint64 = 1

func ShakeHands(conn io.ReadWriter, userAgent string, addr *common.Address) (hs *orcache.Handshake, err error) {
    hs = &orcache.Handshake{}
    if err = orcache.Msg(&orcache.Handshake{Magic, Version, userAgent, addr}).Write(conn); err != nil { return }
    if err = orcache.Msg(hs).Read(conn); err != nil { return }
    if err = orcache.Msg(&orcache.HandshakeAck{}).Write(conn); err != nil { return }
    if err = orcache.Msg(&orcache.HandshakeAck{}).Read(conn); err != nil { return }
    return
}

func MessageListener(conn io.Reader) <-chan butils.Chunk {
    c := make(chan butils.Chunk)
    go func(){
        defer close(c)
        for {
            msg := &orcache.Message{}
            if msg.Read(conn) != nil { return }
            c <- msg.Chunk
        }
    }()
    return c
}

func MessageSender(conn io.WriteCloser) chan<- butils.Chunk {
    c := make(chan butils.Chunk)
    go func(){
        defer conn.Close()
        for {
            chunk, ok := <- c
            if ! ok { return }
            if orcache.Msg(chunk).Write(conn) != nil { return }
        }
    }()
    return c
}
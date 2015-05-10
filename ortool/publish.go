package main
import (
    "fmt"
    "flag"
    "io"
    "orwell/orlib/crypto/card"
    "orwell/orlib/crypto/armor"
    "orwell/orlib/protocol/common"
    "orwell/orlib/client"
    "net"
)

type publishCommand struct{}

func (publishCommand) Usage() string {
    return "ortool publish [--target SERVER] <path>"
}

func (publishCommand) Description() string {
    return `Publishes the given card in the network.

Arguments:
    --target SERVER   Server (addr:port).    Default: 127.0.0.1:1984
`
}

func (publishCommand) Main(args []string) (err error) {

    fs := flag.NewFlagSet("publish", flag.ContinueOnError)
    fTg := fs.String("target", "127.0.0.1:1984", "")
    if fs.Parse(args) != nil { return InvalidUsage }
    if len(fs.Args()) != 1 { return InvalidUsage }

    var r io.Reader
    if r, err = FileOrSTDIN(fs.Arg(0)); err != nil { return }

    c := &card.Card{}
    if err = armor.DecodeFromTo(r, c); err != nil { return }

    var conn net.Conn
    if conn, err = net.Dial("tcp", *fTg); err != nil { return }

    if _, err = client.ShakeHands(conn, "ortool", nil, common.NoPort, nil); err != nil { return }

    var ttl common.TTL
    if ttl, err = client.Publish(conn, c); err != nil { return }

    fmt.Println("Published. TTL =", ttl)
    return
}
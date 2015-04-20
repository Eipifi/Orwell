package main
import (
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/types"
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

    var c *card.Card
    if c, err = card.FromFile(fs.Arg(0)); err != nil { return }

    var ms *orcache.OrcacheMessenger
    if ms, err = orcache.SimpleClient(*fTg); err != nil { return }

    var ttl types.TTL
    if ttl, err = ms.Put(c); err != nil { return }

    fmt.Println("Published. TTL =", ttl)
    return
}
package publish
import (
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/types"
)

const Usage = `usage: ortool publish [--target ip:port] --card path

Publishes the given card in the network.

`
func Main(args []string) {
    var err error

    fs := flag.NewFlagSet("publish", flag.ExitOnError)
    fTg := fs.String("target", "127.0.0.1:1984", "Orcache server address")
    fCd := fs.String("card", "", "Card file path")
    fs.Parse(args)

    var c *card.Card
    if c, err = card.FromFile(*fCd); err != nil {
        fmt.Println("Error:", err)
        return
    }

    var ms *orcache.OrcacheMessenger
    if ms, err = orcache.SimpleClient(*fTg); err != nil {
        fmt.Println("Error:", err)
        return
    }

    var ttl types.TTL
    if ttl, err = ms.Put(c); err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Published. TTL =", ttl)
}
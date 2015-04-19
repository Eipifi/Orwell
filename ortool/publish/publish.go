package publish
import (
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/types"
)

const Usage = `usage: ortool publish [--target ip:port] path

Publishes the given card in the network.

`
func Main(args []string) {

    if len(args) == 0 {
        fmt.Println("Path not provided. See 'ortool help publish' for details.")
        return
    }

    var err error
    fs := flag.NewFlagSet("publish", flag.ExitOnError)
    fTg := fs.String("target", "127.0.0.1:1984", "Orcache server address")
    fs.Parse(args[:len(args)-1])

    var c *card.Card
    if c, err = card.FromFile(args[len(args)-1]); err != nil {
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
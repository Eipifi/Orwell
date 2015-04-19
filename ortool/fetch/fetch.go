package fetch
import (
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/types"
    "orwell/orlib/protocol/orcache"
)

const Usage = `usage: ortool fetch [--from ip:port] <address>

Asks the server for a card.

`
func Main(args []string) {
    if len(args) == 0 {
        fmt.Println("ID not provided.")
        return
    }

    var id *types.ID
    var err error
    if id, err = types.HexToID(args[len(args)-1]); err != nil {
        fmt.Println("Error: ", err)
        return
    }

    args = args[:len(args)-1]

    fs := flag.NewFlagSet("fetch", flag.ExitOnError)
    fSr := fs.String("from", "127.0.0.1:1984", "Orcache server address")
    fs.Parse(args)

    var ms *orcache.OrcacheMessenger
    if ms, err = orcache.SimpleClient(*fSr); err != nil {
        fmt.Println("Error:", err)
        return
    }

    var c *card.Card
    if c, err = ms.Get(id, 0); err != nil {
        fmt.Println(err)
    }

    if c == nil {
        fmt.Println("Card not found on the server.")
    } else {
        fmt.Println(string(c.Payload.MarshalJSON()))
    }
}
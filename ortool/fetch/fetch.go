package fetch
import (
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/types"
    "orwell/orlib/protocol/orcache"
    "encoding/pem"
    "os"
)

const Usage = `usage: ortool fetch [--from ip:port --format json/pem] <address>

Asks the server for a card.

`
func Main(args []string) {
    if len(args) == 0 {
        fmt.Println("ID not provided. See 'ortool help fetch' for details.")
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
    fSrc := fs.String("from", "127.0.0.1:1984", "Orcache server address")
    fFmt := fs.String("format", "json", "Card format [json, pem]")
    fs.Parse(args)

    var ms *orcache.OrcacheMessenger
    if ms, err = orcache.SimpleClient(*fSrc); err != nil {
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
        switch *fFmt {
            case "json":
                fmt.Println(string(c.Payload.MarshalJSON()))
            case "pem":
                block := pem.Block{}
                block.Type = "ORWELL CARD"
                block.Bytes, _ = c.Marshal()
                pem.Encode(os.Stdout, &block)
        }

    }
}
package main
import (
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/types"
    "orwell/orlib/protocol/orcache"
    "encoding/pem"
    "os"
    "errors"
)

type fetchCommand struct {}

func (fetchCommand) Usage() string {
    return "ortool fetch [--from SERVER --format FORMAT] <address>"
}

func (fetchCommand) Description() string {
    return `Resolves the specified address.

For now, accepts only hex address. Todo: implement rest

Arguments:
    --from SERVER    Server (addr:port).    Default: 127.0.0.1:1984
    --format FORMAT  Desired output format. Default: table
    address          Orwell address

Formats:
    table   todo here
    json
    pem
`
}

func (fetchCommand) Main(args []string) (err error) {
    if len(args) == 0 { return InvalidUsage }

    args = args[:len(args)-1]
    fs := flag.NewFlagSet("fetch", flag.ContinueOnError)
    fSrc := fs.String("from", "127.0.0.1:1984", "")
    fFmt := fs.String("format", "json", "")
    if err = fs.Parse(args); err != nil { return InvalidUsage }
    if len(fs.Args()) != 1 { return InvalidUsage }

    var id *types.ID
    if id, err = types.HexToID(fs.Arg(0)); err != nil { return }

    var ms *orcache.OrcacheMessenger
    if ms, err = orcache.SimpleClient(*fSrc); err != nil { return }

    var c *card.Card
    if c, err = ms.Get(id, 0); err != nil { return }

    if c == nil {
        return errors.New("Card not found on the server.")
    } else {
        switch *fFmt {
            case "json":
                fmt.Println(string(c.Payload.MarshalJSON()))
            case "pem":
                block := pem.Block{}
                block.Type = "ORWELL CARD"
                block.Bytes, _ = c.Marshal()
                return pem.Encode(os.Stdout, &block)
        }
    }
    return
}
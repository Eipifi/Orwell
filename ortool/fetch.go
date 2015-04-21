package main
import (
    "flag"
    "orwell/orlib/protocol/types"
    "orwell/orlib/protocol/orcache"
    "os"
    "errors"
    "orwell/orlib/crypto/card"
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
        var b []byte
        switch *fFmt {
            case "json":
                b, err = c.Payload.MarshalJSON()
            case "pem":
                b, err = c.MarshalPEM()
        }
        if err != nil { return }
        _, err = os.Stdout.Write(b)
    }
    return
}
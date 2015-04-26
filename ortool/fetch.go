package main
import (
    "flag"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/conv"
    "orwell/orlib/crypto/card"
    "errors"
    "orwell/orlib/crypto/armor"
    "os"
    "orwell/orlib/butils"
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

    fs := flag.NewFlagSet("fetch", flag.ContinueOnError)
    fSrc := fs.String("from", "127.0.0.1:1984", "")
    fFmt := fs.String("format", "json", "")
    if err = fs.Parse(args); err != nil { return InvalidUsage }
    if len(fs.Args()) != 1 { return InvalidUsage }

    var id *hash.ID
    if id, err = hash.HexToID(fs.Arg(0)); err != nil { return }

    var cv *conv.Conversation
    if cv, err = conv.CreateTCPConversation(*fSrc); err != nil { return }

    if err = cv.DoHandshake("ortool", nil); err != nil { return }

    var c *card.Card
    if c, err = cv.DoGet(id, 0); err != nil { return }

    if c == nil {
        return errors.New("Card not found on the server.")
    } else {
        switch *fFmt {
            case "json":
                var b []byte
                if b, err = c.Payload.WriteJSON(); err != nil { return }
                return butils.WriteFull(os.Stdout, b)
            case "pem":
                return armor.EncodeObjTo(c, armor.TextCard, os.Stdout)
        }
    }
    return
}
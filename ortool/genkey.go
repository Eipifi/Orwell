package main
import (
    "flag"
    "io"
    "orwell/orlib/crypto/sig"
    "orwell/orlib/crypto/armor"
)

type genkeyCommand struct {}

func (genkeyCommand) Usage() string {
    return "ortool genkey [--output OUTPUT]"
}

func (genkeyCommand) Description() string {
    return `Genkey creates a new asymmetric key pair that can be used both
as a card key in orcache and as a wallet in orchain.

Arguments:
    --output OUTPUT         destination (stdout by default)

`
}

func (genkeyCommand) Main(args []string) (err error) {
    fs := flag.NewFlagSet("ortool genkey", flag.ContinueOnError)
    var fOut = fs.String("output", "", "")
    if fs.Parse(args) != nil { return InvalidUsage }
    if len(fs.Args()) != 0 { return InvalidUsage }

    key, err := sig.CreateKey()
    if err != nil { return }

    var out io.Writer
    if out, err = FileOrSTDOUT(*fOut); err != nil { return }

    return armor.EncodeObjTo(key, "PRIVATE KEY", out)
}

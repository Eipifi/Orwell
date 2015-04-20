package main
import (
    "flag"
    "orwell/orlib/sig"
    "encoding/pem"
    "io"
    "errors"
)

type genkeyCommand struct {}

func (genkeyCommand) Usage() string {
    return "ortool genkey [--algorithm ALGORITHM] [--output OUTPUT]"
}

func (genkeyCommand) Description() string {
    return `Genkey creates a new asymmetric key pair that can be used both
as a card key in orcache and as a wallet in orchain.

Available key algorithms:

    ecdsa/p-256/sha256      [default]
    rsassa-pss/sha256

Arguments:
    --algorithm ALGORITHM   key algorithm, specified above
    --output OUTPUT         destination (stdout by default)

`
}

const (
    ECDSA = "ecdsa/p-256/sha256"
)

func (genkeyCommand) Main(args []string) (err error) {
    fs := flag.NewFlagSet("ortool genkey", flag.ContinueOnError)
    var fAlg = fs.String("algorithm", ECDSA, "")
    var fOut = fs.String("output", "", "")
    if fs.Parse(args) != nil { return InvalidUsage }
    if len(fs.Args()) != 0 { return InvalidUsage }

    var key sig.PrvKey

    switch *fAlg {
        case ECDSA:
            key = sig.NewEcdsaPrvKey()
        default:
            return errors.New("Unknown algorithm " + *fAlg)
    }

    var out io.Writer
    if out, err = FileOrSTDOUT(*fOut); err != nil { return }

    block := pem.Block{}
    block.Type = "ORWELL PRIVATE KEY"
    block.Bytes = key.Serialize()
    return pem.Encode(out, &block)
}

package genkey
import (
    "flag"
    "fmt"
    "orwell/orlib/sig"
    "encoding/pem"
    "os"
    "io"
)

const Usage = `usage: ortool genkey [-a algorithm] [-o output]

Genkey creates a new asymmetric key pair that can be used both
as a card key in orcache and as a wallet in orchain.

Available key algorithms:

    ecdsa/p-256/sha256      [default]
    rsassa-pss/sha256

Arguments:
    -algorithm ALGORITHM   key algorithm, specified above
    -output OUTPUT         destination (stdout by default)

`

// TODO: enable short argument forms

func Main(args []string) {
    fs := flag.NewFlagSet("ortool genkey", flag.ExitOnError)

    var fAlg = fs.String("algorithm", "ecdsa/p-256/sha256", "key algorithm")
    var fOut = fs.String("output", "", "destination")

    fs.Parse(args)
    var key sig.PrvKey
    var out io.Writer

    switch *fAlg {
        case "ecdsa/p-256/sha256":
            key = sig.NewEcdsaPrvKey()
        default:
            fmt.Print("Unknown algorithm \"", *fAlg, "\".")
            return
    }

    if *fOut == "" {
        out = os.Stdout
    } else {
        o, e := os.Create(*fOut)
        if e != nil {
            fmt.Println(e.Error())
        } else {
            out = o
        }
    }

    block := pem.Block{}
    block.Type = "ORWELL PRIVATE KEY"
    block.Bytes = key.Serialize()
    pem.Encode(out, &block)
}

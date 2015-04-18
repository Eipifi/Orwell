package gencard
import (
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "orwell/orlib/card"
    "orwell/orlib/sig"
    "encoding/pem"
)

const Usage = `usage: ortool gencard --payload PAYLOAD --key KEY

Gencard creates a new card from a JSON card payload and private key.

Example payload structure:

{
    "version": 42,
    "expires": 123123123,
    "records": [
        {
            "key": "kittens",
            "type": "http",
            "value": "ip=192.168.0.10 ca=cafebabeaa"
        }
    ]
}

`

func Main(args []string) {
    fs := flag.NewFlagSet("gencard", flag.ExitOnError)

    fPld := fs.String("payload", "", "JSON card payload file")
    fKey := fs.String("key", "", "Private key file")

    fs.Parse(args)

    pld, err := ioutil.ReadFile(*fPld)
    if err != nil {
        fmt.Printf("Payload: %s\n", err.Error())
        os.Exit(1)
    }

    keyB, err := ioutil.ReadFile(*fKey)
    if err != nil {
        fmt.Printf("Key: %s\n", err.Error())
        os.Exit(1)
    }

    b, _ := pem.Decode(keyB)

    if b == nil {
        fmt.Printf("Key: failed to strip the PEM armor\n")
        os.Exit(1)
    }

    key, err := sig.ParsePrvKey(b.Bytes)
    if err != nil {
        fmt.Printf("Key: %s\n", err.Error())
        os.Exit(1)
    }

    c, err := card.Create(pld, key)
    if err != nil {
        fmt.Printf("Card: %s\n", err.Error())
        os.Exit(1)
    }

    cb, err := c.Marshal()
    if err != nil {
        fmt.Printf("Card: %s\n", err.Error())
        os.Exit(1)
    }

    block := pem.Block{}
    block.Type = "ORWELL CARD"
    block.Bytes = cb
    pem.Encode(os.Stdout, &block)
}
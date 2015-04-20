package main
import (
    "os"
    "io/ioutil"
    "orwell/orlib/card"
    "orwell/orlib/sig"
    "encoding/pem"
    "errors"
)

type gencardCommand struct{}

func (gencardCommand) Usage() string {
    return "ortool gencard <key> [payload]"
}

func (gencardCommand) Description() string {
    return `Gencard creates a new card from a JSON card payload and private key.

Arguments:
    key      PEM key file
    payload  Path to JSON card file. If empty, STDIN will be used.

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
}

func (gencardCommand) Main(args []string) (err error) {
    if len(args) > 2 || len(args) < 1 { return InvalidUsage }

    var rawKey []byte
    if rawKey, err = ioutil.ReadFile(args[0]); err != nil { return }

    var rawJson []byte
    if rawJson, err = ReadWholeFileOrSTDIN(rs(args, 1)); err != nil { return }

    b, _ := pem.Decode(rawKey)
    if b == nil {
        return errors.New("Key: failed to strip the PEM armor\n")
    }

    var key sig.PrvKey
    if key, err = sig.ParsePrvKey(b.Bytes); err != nil { return }

    var c *card.Card
    if c, err = card.Create(rawJson, key); err != nil { return }

    var cb []byte
    if cb, err = c.Marshal(); err != nil { return }

    block := pem.Block{}
    block.Type = "ORWELL CARD"
    block.Bytes = cb
    return pem.Encode(os.Stdout, &block)
}
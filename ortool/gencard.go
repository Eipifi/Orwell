package main
import (
    "os"
    "orwell/orlib/crypto/sig"
    "io/ioutil"
    "orwell/orlib/crypto/card"
    "orwell/orlib/crypto/armor"
)

type gencardCommand struct{}

func (gencardCommand) Usage() string {
    return "ortool gencard <key> <payload>"
}

func (gencardCommand) Description() string {
    return `Gencard creates a new card from a JSON card payload and private key.

Arguments:
    key      Path to PEM key file
    payload  Path to JSON card file

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
    if len(args) != 2 { return InvalidUsage }

    keyPEM, err := ioutil.ReadFile(args[0])
    if err != nil { return }
    pldJSON, err := ioutil.ReadFile(args[1])
    if err != nil { return }

    key := &sig.PrivateKey{}
    if err = armor.DecodeTo(keyPEM, key); err != nil { return }

    pld := &card.Payload{}
    if err = pld.ReadJSON(pldJSON); err != nil { return }

    card, err := pld.Sign(key)
    if err != nil { return }

    return armor.EncodeObjTo(card, armor.TextCard, os.Stdout)
}
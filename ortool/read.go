package main
import (
    "fmt"
    "io/ioutil"
    "encoding/pem"
    "orwell/orlib/sig"
    "orwell/orlib/card"
    "io"
    "errors"
)

type readCommand struct {}

func (readCommand) Usage() string {
    return "ortool read [path]"
}

func (readCommand) Description() string {
    return `Reads the given file (card or key) and displays relevant info.

Arguments:
    path        File to read. If empty, STDIN will be used
`
}


func (readCommand) Main(args []string) (err error) {
    if len(args) > 1 { return InvalidUsage }

    var f io.Reader
    if f, err = FileOrSTDIN(rs(args, 0)); err != nil { return }

    var data []byte
    if data, err = ioutil.ReadAll(f); err != nil { return }

    b, _ := pem.Decode(data)
    if b == nil { return errors.New("Failed to read PEM file.") }

    key, err := sig.ParsePrvKey(b.Bytes)
    if err == nil {
        fmt.Println("PRIVATE KEY")
        fmt.Printf("ID: %s\n", key.PublicPart().Id())
        return
    }

    c, err := card.Unmarshal(b.Bytes)
    if err == nil {
        fmt.Println("CARD")
        fmt.Printf("ID: %s\n", c.Key.Id())
        fmt.Printf("%s\n", c.Payload.MarshalJSON())
        return
    }
    return errors.New("Failed to parse input.")
}
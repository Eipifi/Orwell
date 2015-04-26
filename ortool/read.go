package main
import (
    "orwell/orlib/crypto/armor"
    "orwell/orlib/crypto/sig"
    "io"
    "fmt"
    "orwell/orlib/crypto/card"
    "errors"
    "os"
    "orwell/orlib/butils"
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

    data, err := armor.DecodeAll(f)
    if err != nil { return }

    k := sig.PrivateKey{}
    if err = k.ReadBytes(data); err == nil {
        fmt.Println("PRIVATE KEY")
        fmt.Println("ID:", k.PublicPart().Id())
        return
    }

    c := card.Card{}
    if err = c.ReadBytes(data); err == nil {
        fmt.Println("CARD")
        fmt.Println("ID:", c.Key.Id())
        var json []byte
        json, err = c.Payload.WriteJSON()
        if err != nil { return }
        return butils.WriteFull(os.Stdout, json)
    } else {
        return err
    }

    return errors.New("PEM content not recognized")
}
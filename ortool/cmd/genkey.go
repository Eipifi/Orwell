package cmd
import (
    "github.com/codegangsta/cli"
    "orwell/lib/crypto/sig"
    "encoding/pem"
    "fmt"
)

func GenKeyAction(c *cli.Context) {
    if err := gen(); err != nil { panic(err) }
}

func gen() error {
    prv, err := sig.Create()
    if err != nil { return err }

    key_contents, err := prv.WriteBytes()
    if err != nil { return err }
    pem_contents := pem.EncodeToMemory(&pem.Block{
        Type: "ORWELL PRIVATE KEY",
        Bytes: key_contents,
    })
    fmt.Println(string(pem_contents))
    return nil
}
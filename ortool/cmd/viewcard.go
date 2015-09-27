package cmd
import (
    "github.com/codegangsta/cli"
    "orwell/lib/protocol/orcache"
    "orwell/lib/crypto/armor"
    "orwell/lib/butils"
    "fmt"
)

func ViewCardAction(ctx *cli.Context) {
    card := orcache.Card{}
    err := armor.ReadFromFile(&butils.BRWrapper{&card}, ctx.Args().First())
    if err != nil {
        fmt.Printf("Error: %v \n", err)
    } else {
        fmt.Printf("ID     : %v \n", card.ID())
        fmt.Printf("Version: %v \n", card.Version)
        fmt.Printf("Expires: %v \n", card.Timestamp)
        fmt.Println("Entries:")
        for _, r := range card.Entries {
            fmt.Printf("  - KEY=%v TYPE=%v VALUE=%v \n", string(r.Key), string(r.Type), string(r.Value))
        }
    }
}
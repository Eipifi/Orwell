package cmd
import (
    "github.com/codegangsta/cli"
    "orwell/lib/fcli"
    "fmt"
    "orwell/lib/crypto/sig"
    "orwell/lib/crypto/armor"
    "orwell/lib/protocol/orcache"
    "orwell/lib/butils"
)

func GenCardAction(ctx *cli.Context) {

    var c orcache.Card
    var key sig.PrvKey


    fmt.Println("Enter the file name of the private key:")
    fsm := fcli.NewFSM("> ")

    fsm.On("init", "$str", func(path string) fcli.Result {
        if err := armor.ReadFromFile(&key, path); err != nil { return err }
        return fcli.Next("main")
    })

    fsm.On("main", "status", func() fcli.Result {
        fmt.Printf("ID     : %v \n", key.PublicPart().ID())
        fmt.Printf("Version: %v \n", c.Version)
        fmt.Printf("Expires: %v \n", c.Timestamp)
        fmt.Println("Entries:")
        for _, r := range c.Entries {
            fmt.Printf("  - KEY=%v TYPE=%v VALUE=%v \n", string(r.Key), string(r.Type), string(r.Value))
        }
        return nil
    })

    fsm.On("main", "version $uint64", func(v uint64) fcli.Result {
        c.Version = v
        return nil
    })

    fsm.On("main", "expires $uint64", func(e uint64) fcli.Result {
        c.Timestamp = e
        return nil
    })

    fsm.On("main", "entry add $str $str $str", func(k, t, v string) fcli.Result {
        rec := orcache.Entry{[]byte(k), []byte(t), []byte(v)}
        c.Entries = append(c.Entries, rec)
        return nil
    })

    fsm.On("main", "commit $str", func(path string) fcli.Result {
        if err := c.Proof.SignHead(&c, &key); err != nil { return err }
        if err := armor.WriteToFile(&butils.BWWrapper{&c}, "ORWELL CARD", path); err != nil { return err }
        return fcli.Exit(nil)
    })

    fsm.On("main", "exit", fcli.ExitHandler)
    fsm.On("main", "x", fcli.ExitHandler)
    fsm.Run("init")
}
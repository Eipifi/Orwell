package command
import (
    "fmt"
    "orwell/lib/db"
    "errors"
)

type ResolveCmd struct{}

func (*ResolveCmd) Name() string {
    return "resolve"
}

func (*ResolveCmd) Run(args []string) error {

    if len(args) != 1 {
        return errors.New("Usage: resolve <domain>")
    }

    domain_name := args[0]

    db.Get().View(func(t *db.Tx){
        domain := t.GetValidRegisteredDomain(domain_name)
        state := t.GetState()
        if domain == nil {
            fmt.Printf("Domain \"%s\" not registered.\n", domain_name)
        } else {
            fmt.Println("Domain:")
            fmt.Println("NAME        : ", domain.Name)
            fmt.Println("OWNER       : ", domain.Owner)
            fmt.Println("VALID UNTIL : #", domain.ValidUntilBlock, " (currently #", state.Length, ")")
        }
    })

    return nil
}
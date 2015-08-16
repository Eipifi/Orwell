package cmd
import (
    "orwell/lib/fcli"
    "orwell/lib/db"
    "fmt"
)

func ResolveHandler(domain_name string) fcli.Result {

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
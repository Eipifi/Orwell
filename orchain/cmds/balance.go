package cmds
import (
    "orwell/lib/db"
    "fmt"
    "orwell/lib/foo"
    "orwell/lib/fcli"
)

func BalanceHandler(id foo.U256) fcli.Result {

    db.Get().View(func(t *db.Tx) {
        fmt.Printf("Balance: %v \n", t.GetBalance(id))
    })

    return nil
}
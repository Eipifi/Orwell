package command
import (
    "fmt"
    "orwell/lib/db"
    "errors"
    "orwell/lib/foo"
)

type BalanceCmd struct{}

func (*BalanceCmd) Name() string {
    return "balance"
}

func (*BalanceCmd) Run(args []string) error {

    if len(args) != 1 {
        return errors.New("Invalid number of arguments")
    }

    wallet_id, err := foo.FromHex(args[0])
    if err != nil { return err }

    db.Get().View(func(t *db.Tx) {
        fmt.Printf("Balance: %v \n", t.GetBalance(wallet_id))
    })

    return nil
}
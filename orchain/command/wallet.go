package command
import (
    "fmt"
    "orwell/lib/wallet"
    "errors"
    "orwell/lib/db"
)

type WalletCmd struct{}

func (*WalletCmd) Name() string { return "wallet" }

func (c *WalletCmd) Run(args []string) error {
    if len(args) == 0 {
        return c.List()
    }
    switch args[0] {
        case "generate": return c.Generate()
        default: return errors.New("Unknown command")
    }
}

func (*WalletCmd) List() error {
    fmt.Println("Wallets:")
    for _, w := range wallet.ListWallets() {
        var balance uint64
        db.Get().View(func(t *db.Tx) {
            balance = t.GetBalance(w.ID())
        })
        fmt.Printf("   * %v (balance: %v) \n", w.ID(), balance)
    }
    return nil
}

func (*WalletCmd) Generate() error {
    w := wallet.Generate()
    if err := w.ExportDefault(); err != nil { return err }
    fmt.Printf("Generated wallet %v \n", w.ID())
    return nil
}
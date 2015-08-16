package cmds
import (
    "orwell/lib/fcli"
    "orwell/lib/wallet"
    "fmt"
    "orwell/lib/db"
)

func WalletGenerateHandler() fcli.Result {
    w := wallet.Generate()
    if err := w.ExportDefault(); err != nil { return err }
    fmt.Printf("Generated wallet %v \n", w.ID())
    return nil
}

func WalletHandler() fcli.Result {
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
package command
import (
    "fmt"
    "orwell/lib/wallet"
    "errors"
    "orwell/lib/config"
    "io/ioutil"
    "strings"
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
    files, err := ioutil.ReadDir(config.Path(config.WALLET_DIR))
    if err != nil { return err }
    fmt.Println("Wallets:")
    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".pem") {
            w, err := wallet.Load(walletPath(f.Name()))
            if err != nil {
                fmt.Printf("wallet %v: %v\n", f.Name(), err)

            } else {
                var balance uint64
                db.Get().View(func(t *db.Tx) {
                    balance = t.GetBalance(w.ID())
                })
                fmt.Printf("   * %v (balance: %v) \n", w.ID(), balance)
            }
        }
    }
    return nil
}

func (*WalletCmd) Generate() error {
    w := wallet.Generate()
    w.Export(walletPath(w.ID().String() + ".pem"), 0700)
    fmt.Printf("Generated wallet %v \n", w.ID())
    return nil
}

func walletPath(name string) string {
    return config.Path(config.WALLET_DIR + "/" + name)
}
package command
import (
    "orwell/lib/wallet"
    "errors"
    "fmt"
    "orwell/lib/cmd"
    "orwell/lib/protocol/orchain"
    "orwell/lib/db"
)

type PayCmd struct{}

func (*PayCmd) Name() string {
    return "pay"
}

func (*PayCmd) Run(args []string) error {

    wallets := wallet.ListWallets()
    if len(wallets) == 0 { return errors.New("No wallet found") }

    fmt.Println("Choose the wallet to send from:")

    for i, w := range wallets {
        fmt.Printf("#%v - %v (balance: %v) \n", i, w.ID(), w.Balance())
    }

    num := cmd.ReadUint64(0, uint64(len(wallets)) - 1)
    wallet := wallets[int(num)]
    bill := orchain.Bill{}
    fmt.Printf("Chosen wallet: %v\n", wallet.ID())
    fmt.Println("Enter the recipient address:")
    bill.Target = cmd.ReadU256()
    fmt.Println("Enter the amount:")
    bill.Value = cmd.ReadUint64(0, wallet.Balance())
    fmt.Println("Enter the fee:")
    fee := cmd.ReadUint64(0, wallet.Balance() - bill.Value)

    txn, err := wallet.CreateTransaction([]orchain.Bill{bill}, fee, orchain.Payload{})

    fmt.Printf("%+v \n", txn)

    if err != nil {
        fmt.Println("Error: ", err)
    } else {
        err = db.Get().UpdateE(func(t *db.Tx) error {
            return t.MaybeStoreUnconfirmedTransaction(txn)
        })

        if err != nil {
            fmt.Println("Transaction was not accepted:", err)
        } else {
            fmt.Println("Transaction was created. Please wait until the transaction is accepted in the blockchain.")
        }
    }

    return nil
}
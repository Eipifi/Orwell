package command
import (
    "fmt"
    "orwell/lib/db"
    "orwell/lib/protocol/orchain"
    "orwell/lib/wallet"
    "errors"
    "orwell/lib/cmd"
)

type SendCmd struct{}
func (*SendCmd) Name() string { return "send" }

/*
    Send tool:
        - Choose sender wallet
        - Display balance in loop
        - Manage outputs (send xyz 123)
        - fee 123
        - ticket name owner valid_until
        - announce name owner valid_until
        - label label
        - transfer coming soon
        - cancel
        - commit
*/

func (*SendCmd) Run(args []string) error {
    wallet := pickWallet()
    if wallet == nil { return errors.New("No wallet found") }

    txn := &orchain.Transaction{}
    var sum_input uint64 = 0
    var fee uint64 = 0
    db.Get().View(func(t *db.Tx) {
        txn.Inputs = t.GetUnspentBillsByWallet(wallet.ID())
        for _, inp := range txn.Inputs {
            sum_input += t.GetBill(&inp).Value
        }
    })
    for {
        // display state
        fmt.Println("----------------------------------------------------------------")
        fmt.Println("Sending from ", wallet.ID())
        fmt.Println("Outputs:")
        for _, out := range txn.Outputs {
            fmt.Printf("   recipient: %v amount: %v \n", out.Target, out.Value)
        }
        fmt.Println("Spent in total:", txn.TotalOutput())
        fmt.Println("Fee:", fee)
        fmt.Println("Balance after txn:", sum_input - txn.TotalOutput() - fee)
        fmt.Println("")

        command := cmd.ReadString()
        fmt.Println(command)
    }



    return nil
}

func pickWallet() *wallet.Wallet {
    wallets := wallet.ListWallets()
    if len(wallets) == 0 { return nil }
    fmt.Println("Choose the wallet to send from:")
    for i, w := range wallets {
        fmt.Printf("#%v - %v (balance: %v) \n", i, w.ID(), w.Balance())
    }
    num := cmd.ReadUint64(0, uint64(len(wallets)) - 1)
    return &wallets[int(num)]
}
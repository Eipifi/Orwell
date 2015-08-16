package cmds
import (
    "orwell/lib/fcli"
    "orwell/lib/wallet"
    "fmt"
    "errors"
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/db"
)

func SendHandler() fcli.Result {

    wallets := wallet.ListWallets()
    var w *wallet.Wallet
    txn := orchain.Transaction{}
    var fee uint64 = 0
    var sum_input uint64 = 0

    print_txn := func() {
        spent := txn.TotalOutput() + fee
        remaining := sum_input - spent
        fmt.Println("---------------------------------------------------------")
        fmt.Println("Sending from   :", w.ID())
        fmt.Println("Fee            :", fee)
        fmt.Println("Spent in total :", spent)
        fmt.Println("Remaining      :", remaining)
        fmt.Println("Payload        :", txn.Payload.String())
        fmt.Println("Outputs :")
        for _, out := range txn.Outputs {
            fmt.Printf(" - %v amount: %v \n", out.Target, out.Value)
        }
    }

    fsm := fcli.NewFSM("send> ")

    // First, select the wallet to pay with
    fsm.On("init", "", func() fcli.Result {
        if len(wallets) == 0 { return fcli.Exit(errors.New("Please generate a wallet first")) }
        fmt.Println("Choose the wallet to send from:")
        for i, w := range wallets {
            fmt.Printf("#%v - %v (balance: %v) \n", i, w.ID(), w.Balance())
        }
        return fcli.Next("choose_wallet")
    })

    // After the wallet is selected, configure options
    fsm.On("choose_wallet", "$uint64", func(id uint64) fcli.Result {
        w = &wallets[id]
        db.Get().View(func(t *db.Tx) {
            txn.Inputs = t.GetUnspentBillsByWallet(w.ID())
            for _, inp := range txn.Inputs {
                sum_input += t.GetBill(&inp).Value
            }
        })
        return fcli.Next("main")
    })

    // Add recipient
    fsm.On("main", "pay $U256 $uint64", func(recipient foo.U256, amount uint64) fcli.Result {
        txn.Outputs = append(txn.Outputs, orchain.Bill{Target: recipient, Value: amount})
        return nil
    })

    // Change fee
    fsm.On("main", "fee $uint64", func(new_fee uint64) fcli.Result {
        fee = new_fee
        return nil
    })

    fsm.On("main", "status", func() fcli.Result {
        print_txn()
        return nil
    })

    fsm.On("main", "label $str", func(label string) fcli.Result {
        txn.Payload = orchain.PayloadLabelString(label)
        return nil
    })

    fsm.On("main", "announce $str $U256 $uint64", func(name string, owner foo.U256, valid_until uint64) fcli.Result {
        txn.Payload = orchain.PayloadDomain(orchain.Domain{name, owner, valid_until})
        return nil
    })

    fsm.On("main", "ticket $str $U256 $uint64", func(name string, owner foo.U256, valid_until uint64) fcli.Result {
        domain := orchain.Domain{name, owner, valid_until}
        txn.Payload = orchain.PayloadTicket(domain.ID())
        return nil
    })

    fsm.On("main", "transfer $str $U256 $uint64", func(name string, owner foo.U256, valid_until uint64) fcli.Result {
        domain := orchain.Domain{name, owner, valid_until}
        return db.Get().ViewE(func(t *db.Tx) error {
            registered_domain := t.GetRegisteredDomain(name)
            if registered_domain == nil { return errors.New("Can not transfer domain - domain not registered")}
            for _, w := range wallets {
                if w.ID() == registered_domain.Owner {
                    transfer := orchain.Transfer{}
                    transfer.Domain = domain
                    transfer.Proof.Sign(&domain, w.PrvKey())
                    txn.Payload = orchain.PayloadTransfer(transfer)
                    return nil
                }
            }
            return errors.New("No wallet matches the current domain owner")
        })
    })

    fsm.On("main", "prepare", func() fcli.Result {
        remaining := sum_input - txn.TotalOutput() - fee
        txn.Outputs = append(txn.Outputs, orchain.Bill{w.ID(), remaining})
        err := txn.Sign(w.PrvKey())
        if err != nil { return fcli.Exit(err) }
        fmt.Println("Prepared transaction:")
        print_txn()
        fmt.Println("Type \"commit\" to send transaction or \"cancel\" to stop the payment (last chance to do so!).")
        return fcli.Next("prepared")
    })

    fsm.On("prepared", "cancel", func() fcli.Result {
        fmt.Println("Transaction cancelled.")
        return fcli.Exit(nil)
    })

    fsm.On("prepared", "commit", func() fcli.Result {
        err := db.Get().UpdateE(func(t *db.Tx) error {
            return t.MaybeStoreUnconfirmedTransaction(&txn)
        })
        if err == nil {
            fmt.Println("Transaction sent.")
        }
        return fcli.Exit(err)
    })

    fsm.On("main", "exit", fcli.ExitHandler)
    fsm.On("main", "x", fcli.ExitHandler)

    return fsm.Run("init")
}
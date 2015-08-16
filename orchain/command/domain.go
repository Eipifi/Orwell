package command
import (
    "fmt"
    "orwell/lib/protocol/orchain"
    "orwell/lib/cmd"
    "orwell/lib/foo"
    "orwell/lib/wallet"
    "errors"
    "orwell/lib/db"
)

type DomainCmd struct{}

func (*DomainCmd) Name() string {
    return "domain"
}

func (*DomainCmd) Run(args []string) error {

    wallets := wallet.ListWallets()
    if len(wallets) == 0 { return errors.New("No wallet found") }

    fmt.Println("This Wizard® will take you through the domain registration process.")
    fmt.Println("Do you want to register a new domain or transfer the existing one?")

    fmt.Println("1 - register a new domain")
    fmt.Println("2 - transfer an existing one")

    fmt.Println("Registering a new domain.")
    fmt.Println("Enter the domain name:")

    domain := &orchain.Domain{}
    domain.Name = cmd.ReadString()

    fmt.Println("Enter the domain owner:")
    domain.Owner = cmd.ReadU256()

    fmt.Println("Domain should be valid until which block?")
    domain.ValidUntilBlock = cmd.ReadUint64(0, foo.U64_MAX)
    ticket := domain.ID()

    fmt.Println("Ok, this is the domain we want to register:")
    fmt.Printf("%+v\n", domain)
    fmt.Printf("Ticket: %v \n", ticket)
    fmt.Println("This ticket will be published right now.")

    for i, w := range wallets {
        fmt.Printf("#%v - %v (balance: %v) \n", i, w.ID(), w.Balance())
    }

    fmt.Println("Choose the wallet to send from:")
    num := cmd.ReadUint64(0, uint64(len(wallets)) - 1)
    wallet := wallets[int(num)]

    fmt.Println("Enter the fee:")
    fee := cmd.ReadUint64(0, foo.U64_MAX)
    txn, err := wallet.CreateTransaction([]orchain.Bill{}, fee, orchain.Payload{Ticket: &ticket})

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
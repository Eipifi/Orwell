package cmd
import (
    "orwell/lib/fcli"
    "fmt"
    "orwell/lib/db"
    "errors"
)


func BlockByNumHandler(num uint64) fcli.Result {

    return db.Get().ViewE(func(t *db.Tx) error {
        bid := t.GetIDByNum(num)
        if bid == nil { return errors.New("Block not found") }
        b := t.GetBlock(*bid)

        fmt.Printf("READ %+v \n", b)

        fmt.Printf("Block       : %v (#%v)\n", *bid, num)
        fmt.Printf("Previous    : %v\n", b.Header.Previous)
        fmt.Printf("Difficulty  : %v \n", b.Header.Difficulty)
        fmt.Printf("Merkle root : %v\n", b.Header.MerkleRoot)
        fmt.Printf("Nonce       : %v\n", b.Header.Nonce)
        fmt.Printf("Timestamp   : %v\n", b.Header.Timestamp)
        fmt.Printf("Transactions:\n")

        for n, tx := range b.Transactions {

            var inputs uint64
            var outputs uint64 = tx.TotalOutput()
            for _, inp := range tx.Inputs {
                inputs += t.GetBill(&inp).Value
            }

            fmt.Printf("   -------- Transaction #%v -------- \n", n)
            fmt.Printf("   ID     : %v\n", tx.ID())
            fmt.Printf("   Input  : %v\n", inputs)
            fmt.Printf("   Output : %v\n", outputs)
            if outputs > inputs {
                fmt.Printf("   Fee    : 0 (coinbase) \n")
            } else {
                fmt.Printf("   Fee    : %v\n", inputs-outputs)
            }
            fmt.Printf("   Inputs:\n")
            for _, inp := range tx.Inputs {
                bill := t.GetBill(&inp)
                fmt.Printf("    [bill=%v/#%v] from=%v value=%v\n", inp.Txn, inp.Index, bill.Target, bill.Value)
            }
            fmt.Printf("   Outputs:\n")
            for n2, out := range tx.Outputs {
                fmt.Printf("    [bill=%v/#%v] to=%v value=%v\n", tx.ID(), n2, out.Target, out.Value)
            }
            fmt.Printf("   Payload:\n")
            if tx.Payload.Label != nil {
                fmt.Printf("    Label \"%s\"\n", string(*tx.Payload.Label))
            } else if tx.Payload.Ticket != nil {
                fmt.Printf("    Ticket %v\n", *tx.Payload.Ticket)
            } else if tx.Payload.Domain != nil {
                fmt.Printf("    Domain name=%s owner=%v expires=%v (ticket=%v)\n", tx.Payload.Domain.Name, tx.Payload.Domain.Owner, tx.Payload.Domain.ValidUntilBlock, tx.Payload.Domain.ID())
            } else if tx.Payload.Transfer != nil {
                fmt.Printf("    Transfer name=%s owner=%v expires=%v\n", tx.Payload.Transfer.Domain.Name, tx.Payload.Transfer.Domain.Owner, tx.Payload.Transfer.Domain.ValidUntilBlock)
            } else {
                fmt.Printf("    none\n")
            }
            fmt.Printf("\n")
        }

        fmt.Printf("Domains :\n")
        for _, d := range b.Domains {
            fmt.Printf("  name=%s owner=%v expires=%v\n", d.Name, d.Owner, d.ValidUntilBlock)
        }

        return nil
    })
}
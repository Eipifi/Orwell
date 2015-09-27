package db
import (
    "orwell/lib/protocol/orchain"
    "fmt"
)

func (t *Tx) PushBlock(block *orchain.Block) error {
    bid := block.Header.ID()
    t.session = &bid
    if err := t.ValidateNewBlock(block); err != nil { return err }
    s := t.GetState()
    t.PutBlock(block, s.Length)
    for _, domain := range block.Domains {
        t.RegisterDomain(&domain)
    }
    s.Length += 1
    s.Head = bid
    s.Work = s.Work.Plus(block.Header.Difficulty)
    t.PutState(s)
    t.RefreshUnconfirmedTransactions()
    fmt.Printf("PUSHED %+v \n", block)
    return nil
}

func (t *Tx) PopBlock() {
    s := t.GetState()
    if s.Length == 1 { panic("Can not remove the genesis block") }
    t.Rollback(s.Head)
    t.RefreshUnconfirmedTransactions()
}
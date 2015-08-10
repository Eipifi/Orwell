package db
import (
    "orwell/lib/protocol/orchain"
)

func (t *Tx) PushBlock(block *orchain.Block) error {
    bid := block.Header.ID()
    t.session = &bid
    if err := t.ValidateNewBlock(block); err != nil { return err }
    s := t.GetState()
    t.PutBlock(block, s.Length)
    s.Length += 1
    s.Head = bid
    s.Work = s.Work.Plus(block.Header.Difficulty)
    t.PutState(s)
    t.RefreshUnconfirmedTransactions()
    return nil
}

func (t *Tx) PopBlock() {
    s := t.GetState()
    if s.Length == 1 { panic("Can not remove the genesis block") }
    t.Rollback(s.Head)
    t.RefreshUnconfirmedTransactions()
}
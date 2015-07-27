package db
import (
    "orwell/lib/protocol/orchain"
)

func (t *Tx) PushBlock(block *orchain.Block) error {
    if err := t.VerifyNextBlock(block); err != nil { return err }
    s := t.GetState()
    t.PutBlock(block, s.Length)
    s.Length += 1
    s.Head = block.Header.ID()
    s.Work = s.Work.Plus(block.Header.Difficulty)
    t.PutState(s)
    return nil
}

func (t *Tx) PopBlock() {
    s := t.GetState()
    if s.Length == 1 { panic("Can not remove the genesis block") }
    h := t.GetHeaderByID(s.Head)
    t.DelBlock(s.Head)
    s.Length -= 1
    s.Head = h.Previous
    s.Work = s.Work.Minus(h.Difficulty)
    t.PutState(s)
}
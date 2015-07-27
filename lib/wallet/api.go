package wallet
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/sig"
    "orwell/lib/utils"
    "orwell/lib/db"
)

type Wallet struct {
    key sig.PrvKey
}

func (w *Wallet) ID() foo.U256 {
    id, err := w.key.PublicPart().ID()
    utils.Ensure(err)
    return id
}

func (w *Wallet) Balance() (sum uint64) {
    db.Get().View(func(t *db.Tx) {
        for _, inp := range t.GetUnspentBillsByUser(w.ID()) {
            sum += t.GetBill(&inp).Value
        }
    })
    return
}
package wallet
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/sig"
    "orwell/lib/utils"
    "orwell/lib/db"
)

type Wallet struct {
    PrvKey sig.PrvKey
}

func (w *Wallet) ID() foo.U256 {
    id, err := w.PrvKey.PublicPart().TryID()
    utils.Ensure(err)
    return id
}

func (w *Wallet) Balance() (balance uint64) {
    db.Get().View(func(t *db.Tx) {
        balance = t.GetBalance(w.ID())
    })
    return
}

func Generate() (w *Wallet) {
    w = &Wallet{}
    key, err := sig.Create()
    utils.Ensure(err)
    w.PrvKey = *key
    return
}
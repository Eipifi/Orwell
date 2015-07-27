package wallet
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/sig"
    "orwell/lib/utils"
)

type Wallet struct {
    key sig.PrvKey
}

func (w *Wallet) ID() foo.U256 {
    id, err := w.key.PublicPart().ID()
    utils.Ensure(err)
    return id
}
package wallet
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/sig"
    "orwell/lib/utils"
    "orwell/lib/db"
    "orwell/lib/protocol/orchain"
    "errors"
)

type Wallet struct {
    key sig.PrvKey
}

func (w *Wallet) ID() foo.U256 {
    id, err := w.key.PublicPart().TryID()
    utils.Ensure(err)
    return id
}

func (w *Wallet) CreateTransaction(bills []orchain.Bill, fee uint64, pld orchain.Payload) (txn *orchain.Transaction, err error) {
    // TODO: check for overflows
    id := w.ID()
    var sum_input, sum_output uint64
    txn = &orchain.Transaction{}
    txn.Outputs = bills
    txn.Payload = pld
    sum_output = txn.TotalOutput()
    db.Get().View(func(t *db.Tx) {
        txn.Inputs = t.GetUnspentBillsByWallet(id)
        for _, inp := range txn.Inputs {
            sum_input += t.GetBill(&inp).Value
        }
    })
    if sum_output + fee > sum_input {
        return nil, errors.New("Not enough funds in wallet")
    }
    rest := sum_input - sum_output - fee
    if rest > 0 {
        txn.Outputs = append(txn.Outputs, orchain.Bill{id, rest})
    }
    if err = txn.Sign(&w.key); err != nil { return }
    return
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
    w.key = *key
    return
}
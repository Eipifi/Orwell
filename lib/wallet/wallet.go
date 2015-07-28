package wallet
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/sig"
    "orwell/lib/utils"
    "orwell/lib/db"
    "orwell/lib/protocol/orchain"
    "errors"
    "os"
    "io/ioutil"
    "encoding/pem"
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
        for _, inp := range t.GetUnspentBillsByWallet(w.ID()) {
            sum += t.GetBill(&inp).Value
        }
    })
    return
}

func (w *Wallet) CreateTransaction(bills []orchain.Bill, fee uint64, label string) (txn *orchain.Transaction, err error) {
    // TODO: check for overflows
    id := w.ID()
    var sum_input, sum_output uint64
    txn = &orchain.Transaction{}
    txn.Label = label
    txn.Outputs = bills
    for _, out := range txn.Outputs {
        sum_output += out.Value
    }
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
    return
}

func Load(file string) (*Wallet, error) {
    w := &Wallet{}
    file_contents, err := ioutil.ReadFile(file)
    if err != nil { return nil, err }
    block, _ := pem.Decode(file_contents)
    if block == nil { return nil, errors.New("Failed to parse PEM block") }
    if err = w.key.ReadBytes(block.Bytes); err != nil { return nil, err }
    return w, nil
}

func (w *Wallet) Export(file string, perm os.FileMode) error {
    key_contents, err := w.key.WriteBytes()
    if err != nil { return err }
    pem_contents := pem.EncodeToMemory(&pem.Block{
        Type: "ORWELL PRIVATE KEY",
        Bytes: key_contents,
    })
    return ioutil.WriteFile(file, pem_contents, perm)
}

func Generate() (w *Wallet) {
    w = &Wallet{}
    key, err := sig.Create()
    utils.Ensure(err)
    w.key = *key
    return
}
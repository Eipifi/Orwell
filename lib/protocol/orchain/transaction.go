package orchain
import (
    "io"
    "errors"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
    "orwell/lib/foo"
    "orwell/lib/utils"
    "orwell/lib/crypto/sig"
)

const TXN_MAX_OUT uint64 = 128
const TXN_MAX_IN uint64 = 65536 // TODO: do something about it

type Transaction struct {
    Payload Payload
    Inputs []BillNumber
    Outputs []Bill
    Proof *sig.Proof // optional - coinbase transactions do not have proofs
}

func (t *Transaction) ReadHead(r io.Reader) (err error) {
    if err = t.Payload.Read(r); err != nil { return }
    if err = butils.ReadSlice(r, TXN_MAX_IN, &t.Inputs); err != nil { return }
    if err = butils.ReadSlice(r, TXN_MAX_OUT, &t.Outputs); err != nil { return }
    return
}

func (t *Transaction) WriteHead(w io.Writer) (err error) {
    if err = t.Payload.Write(w); err != nil { return }
    if err = butils.WriteSlice(w, TXN_MAX_IN, t.Inputs); err != nil { return }
    if err = butils.WriteSlice(w, TXN_MAX_OUT, t.Outputs); err != nil { return }
    return
}

func (t *Transaction) Read(r io.Reader) (err error) {
    if err = t.ReadHead(r); err != nil { return }
    var flag byte
    proof := &sig.Proof{}
    if flag, err = butils.ReadOptional(r, proof); err != nil { return }
    if flag != 0x00 { t.Proof = proof }
    return
}

func (t *Transaction) Write(w io.Writer) (err error) {
    if err = t.WriteHead(w); err != nil { return }
    return butils.WriteOptional(w, t.Proof)
}

func (t *Transaction) TryID() (foo.U256, error) {
    return hash.HashOf(t)
}

func (t *Transaction) ID() foo.U256 {
    id, err := t.TryID()
    utils.Ensure(err)
    return id
}

func (t *Transaction) TotalOutput() (sum uint64) {
    for _, out := range t.Outputs {
        sum += out.Value
    }
    return
}

// This method only verifies if the signature correctly signs the transaction head.
// To ensure the correctness of a transaction, you also need to check if the public key matches the transaction inputs.
func (t *Transaction) Verify() (err error) {
    if t.Proof == nil { return errors.New("Proof missing (valid only for generation transaction)") }
    return t.Proof.CheckHead(t)
}

func (t *Transaction) Sign(key *sig.PrvKey) (err error) {
    t.Proof = &sig.Proof{}
    return t.Proof.SignHead(t, key)
}

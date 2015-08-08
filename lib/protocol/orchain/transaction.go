package orchain
import (
    "io"
    "errors"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
    "bytes"
    "orwell/lib/foo"
    "orwell/lib/utils"
    "orwell/lib/crypto/sig"
)

const TXN_MAX_OUT uint64 = 128
const TXN_MAX_IN uint64 = 65536 // TODO: do something about it
const LABEL_MAX_LENGTH uint64 = 64
var ErrArrayTooLarge = errors.New("Array too large")

type Transaction struct {
    Label string
    Inputs []BillNumber
    Outputs []Bill
    Proof *Proof // optional - coinbase transactions do not have proofs
}

func (t *Transaction) ReadHead(r io.Reader) (err error) {
    if t.Label, err = butils.ReadString(r, LABEL_MAX_LENGTH); err != nil { return }
    if err = butils.ReadSlice(r, TXN_MAX_IN, &t.Inputs); err != nil { return }
    if err = butils.ReadSlice(r, TXN_MAX_OUT, &t.Outputs); err != nil { return }
    return
}

func (t *Transaction) WriteHead(w io.Writer) (err error) {
    if err = butils.WriteString(w, t.Label, LABEL_MAX_LENGTH); err != nil { return }
    if err = butils.WriteSlice(w, TXN_MAX_IN, t.Inputs); err != nil { return }
    if err = butils.WriteSlice(w, TXN_MAX_OUT, t.Outputs); err != nil { return }
    return
}

func (t *Transaction) Read(r io.Reader) (err error) {
    if err = t.ReadHead(r); err != nil { return }
    var flag byte
    proof := &Proof{}
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
    buf := &bytes.Buffer{}
    if err = t.WriteHead(buf); err != nil { return }
    return t.Proof.Check(buf.Bytes())
}

func (t *Transaction) Sign(key *sig.PrvKey) (err error) {
    p := &Proof{}
    pk := key.PublicPart()
    p.PublicKey = *pk
    buf := &bytes.Buffer{}
    err = t.WriteHead(buf)
    if err != nil { return }
    sg, err := key.Sign(buf.Bytes())
    if err != nil { return }
    p.Signature = *sg
    t.Proof = p
    return
}

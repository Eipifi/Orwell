package orchain
import (
    "io"
    "errors"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
    "bytes"
    "orwell/lib/foo"
)

const TXN_MAX_OUT uint64 = 128
const TXN_MAX_IN uint64 = 128
const LABEL_MAX_LENGTH uint64 = 64
var ErrArrayTooLarge = errors.New("Array too large")

type Transaction struct {
    Label string
    Inputs []BillNumber
    Outputs []Bill
    Proof *Proof
}

func (t *Transaction) ReadHead(r io.Reader) (err error) {
    if t.Label, err = butils.ReadString(r, LABEL_MAX_LENGTH); err != nil { return }
    var num uint64
    if num, err = butils.ReadVarUint(r); err != nil { return }
    if num > TXN_MAX_IN { return ErrArrayTooLarge }
    t.Inputs = make([]BillNumber, num)
    for i := 0; i < int(num); i += 1 {
        if err = t.Inputs[i].Read(r); err != nil { return }
    }
    if num, err = butils.ReadVarUint(r); err != nil { return }
    if num > TXN_MAX_OUT { return ErrArrayTooLarge }
    t.Outputs = make([]Bill, num)
    for i := 0; i < int(num); i += 1 {
        if err = t.Outputs[i].Read(r); err != nil { return }
    }
    return
}

func (t *Transaction) WriteHead(w io.Writer) (err error) {
    if err = butils.WriteString(w, t.Label, LABEL_MAX_LENGTH); err != nil { return }
    num := uint64(len(t.Inputs))
    if num > TXN_MAX_IN { return ErrArrayTooLarge }
    if err = butils.WriteVarUint(w, num); err != nil { return }
    for i := 0; i < int(num); i += 1 {
        if err = t.Inputs[i].Write(w); err != nil { return }
    }
    num = uint64(len(t.Outputs))
    if num > TXN_MAX_OUT { return ErrArrayTooLarge }
    if err = butils.WriteVarUint(w, num); err != nil { return }
    for i := 0; i < int(num); i += 1 {
        if err = t.Outputs[i].Write(w); err != nil { return }
    }
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

func (t *Transaction) ID() (foo.U256, error) {
    return hash.HashOf(t)
}

// This method only verifies if the signature correctly signs the transaction head.
// To ensure the correctness of a transaction, you also need to check if the public key matches the transaction inputs.
func (t *Transaction) Verify() (err error) {
    if t.Proof == nil { return errors.New("Proof missing (valid only for generation transaction)") }
    buf := &bytes.Buffer{}
    if err = t.WriteHead(buf); err != nil { return }
    return t.Proof.Check(buf.Bytes())
}

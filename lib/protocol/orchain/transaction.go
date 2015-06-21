package orchain
import (
    "io"
    "errors"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
    "bytes"
)

const TXN_MAX_OUT uint64 = 128
const TXN_MAX_IN uint64 = 128
var ErrArrayTooLarge = errors.New("Array too large")

type Transaction struct {
    TimeLock uint64                     // last acceptable block number
    //NameHash *id.ID                   // hash of the DomainData structure
    Outputs []Bill                    // list of outputs
    Inputs []BillNumber                      // list of inputs
    Proofs []Proof                      // list of public keys and signatures corresponding to inputs
}

func (t *Transaction) Read(r io.Reader) (err error) {
    if t.TimeLock, err = butils.ReadUint64(r); err != nil { return }
    var num uint64

    if num, err = butils.ReadVarUint(r); err != nil { return }
    if num > TXN_MAX_OUT { return ErrArrayTooLarge }
    t.Outputs = make([]Bill, num)
    for i := 0; i < int(num); i += 1 {
        if err = t.Outputs[i].Read(r); err != nil { return }
    }

    if num, err = butils.ReadVarUint(r); err != nil { return }
    if num > TXN_MAX_IN { return ErrArrayTooLarge }
    t.Inputs = make([]BillNumber, num)
    for i := 0; i < int(num); i += 1 {
        if err = t.Inputs[i].Read(r); err != nil { return }
    }

    t.Proofs = make([]Proof, num)
    for i := 0; i < int(num); i += 1 {
        if err = t.Proofs[i].Read(r); err != nil { return }
    }
    return nil
}

func (t *Transaction) WriteHead(w io.Writer) (err error) {
    if err = butils.WriteUint64(w, t.TimeLock); err != nil { return }

    num := uint64(len(t.Outputs))
    if num > TXN_MAX_OUT { return ErrArrayTooLarge }
    if err = butils.WriteVarUint(w, num); err != nil { return }
    for i := 0; i < int(num); i += 1 {
        if err = t.Outputs[i].Write(w); err != nil { return }
    }

    num = uint64(len(t.Inputs))
    if num > TXN_MAX_IN { return ErrArrayTooLarge }
    if err = butils.WriteVarUint(w, num); err != nil { return }
    for i := 0; i < int(num); i += 1 {
        if err = t.Inputs[i].Write(w); err != nil { return }
    }

    return nil
}

func (t *Transaction) Write(w io.Writer) (err error) {
    if err = t.WriteHead(w); err != nil { return }
    for i := 0; i < len(t.Proofs); i += 1 {
        if err = t.Proofs[i].Write(w); err != nil { return }
    }
    return nil
}

func (t *Transaction) ID() (butils.Uint256, error) {
    return hash.HashOf(t)
}

// This method only verifies if the signatures correctly sign the transaction head.
// To ensure the correctness of a transaction, you also need to check if the public keys match the transaction inputs.
func (t *Transaction) VerifySignatures() (err error) {
    if len(t.Inputs) != len(t.Proofs) {
        return errors.New("Mismatch between the number of inputs and proofs")
    }
    buf := &bytes.Buffer{}
    if err = t.WriteHead(buf); err != nil { return }
    head := hash.Hash(buf.Bytes())
    for i := 0; i < len(t.Inputs); i += 1 {
        if err = t.Proofs[i].Check(head); err != nil { return }
    }
    return nil
}

package orchain
import (
    "io"
    "orwell/lib/butils"
    "orwell/lib/foo"
)

type BillNumber struct {
    Txn foo.U256                // transaction being referenced
    Index uint64                // output number in the referenced transaction
}

func (i *BillNumber) Read(r io.Reader) (err error) {
    if err = i.Txn.Read(r); err != nil { return }
    if i.Index, err = butils.ReadVarUint(r); err != nil { return }
    return
}

func (i *BillNumber) Write(w io.Writer) (err error) {
    if err = i.Txn.Write(w); err != nil { return }
    if err = butils.WriteVarUint(w, i.Index); err != nil { return }
    return
}
package orchain
import (
    "io"
    "orwell/lib/butils"
    "orwell/lib/foo"
)

type Bill struct {
    Target foo.U256       // hash of the new owner's public key
    Value uint64                // amount of transferred money
}

func (o *Bill) Read(r io.Reader) (err error) {
    if err = o.Target.Read(r); err != nil { return }
    if o.Value, err = butils.ReadVarUint(r); err != nil { return }
    return
}

func (o *Bill) Write(w io.Writer) (err error) {
    if err = o.Target.Write(w); err != nil { return }
    if err = butils.WriteVarUint(w, o.Value); err != nil { return }
    return
}
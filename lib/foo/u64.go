package foo
import (
    "io"
    "orwell/lib/butils"
)

type U64 uint64

func (u *U64) Read(r io.Reader) error {
    val, err := butils.ReadUint64(r)
    if err != nil { return err }
    *u = U64(val)
    return nil
}

func (u *U64) Write(w io.Writer) error {
    return butils.WriteUint64(w, uint64(*u))
}

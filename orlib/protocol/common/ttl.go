package common
import (
    "io"
    "orwell/orlib/butils"
)

const MaxTTLValue = TTL(255)
type TTL uint8

func (t *TTL) Read(r io.Reader) error {
    val, err := butils.ReadUint8(r)
    if err != nil { return err }
    *t = TTL(val)
    return nil
}

func (t *TTL) Write(w io.Writer) error {
    return butils.WriteUint8(w, uint8(*t))
}
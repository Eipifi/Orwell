package obp
import (
    "orwell/lib/butils"
    "io"
)

const FRAME_MAX_LENGTH uint64 = 1024 * 1024

type Frame struct {
    Context uint64
    Payload []byte
}

func (f *Frame) Read(r io.Reader) (err error) {
    if f.Context, err = butils.ReadVarUint(r); err != nil { return }
    if f.Payload, err = butils.ReadVarBytes(r, FRAME_MAX_LENGTH); err != nil { return }
    return nil
}

func (f *Frame) Write(w io.Writer) (err error) {
    if err = butils.WriteVarUint(w, f.Context); err != nil { return }
    if err = butils.WriteVarBytes(w, f.Payload, FRAME_MAX_LENGTH); err != nil { return }
    return nil
}
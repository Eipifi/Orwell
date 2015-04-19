package types
import "orwell/orlib/comm"

const MaxFrameLength = 65535

type Frame struct {
    Command uint64
    Payload []byte
}

func (f *Frame) Read(r *comm.Reader) (err error) {
    if f.Command, err = r.ReadVaruint(); err != nil { return }
    if f.Payload, err = r.ReadVarBytes(); err != nil { return }
    return
}

func (f *Frame) Write(w *comm.Writer) {
    w.WriteVaruint(f.Command)
    w.WriteVarBytes(f.Payload)
}
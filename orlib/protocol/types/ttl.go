package types
import "orwell/orlib/comm"

const MaxTTLValue = TTL(255)
type TTL uint8

func (t *TTL) Read(r *comm.Reader) (err error) {
    var v uint8
    v, err = r.ReadUint8()
    *t = TTL(v)
    return
}

func (t *TTL) Write(w *comm.Writer) {
    w.WriteUint8(uint8(*t))
}
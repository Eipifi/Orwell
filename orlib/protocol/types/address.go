package types
import "orwell/orlib/comm"

type Address struct {
    IP [16]byte
    Port uint16
    Nonce uint64
}

func (a *Address) Read(r *comm.Reader) (err error) {
    if err = r.ReadTo(a.IP[:]); err != nil { return }
    if a.Port, err = r.ReadUint16(); err != nil { return }
    if a.Nonce, err = r.ReadUint64(); err != nil { return }
    return
}

func (a *Address) Write(w *comm.Writer) {
    w.Write(a.IP[:]) // 16 bytes
    w.WriteUint16(a.Port)
    w.WriteUint64(a.Nonce)
}

func (a *Address) Id() *ID {
    var b [24]byte
    copy(b[0:16], a.IP[:])
    comm.ByteOrder.PutUint64(b[16:24], a.Nonce)
    return Hash(b[:])
}
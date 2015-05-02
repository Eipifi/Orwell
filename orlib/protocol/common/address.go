package common
import (
    "io"
    "orwell/orlib/butils"
    "net"
    "strconv"
    "errors"
    "orwell/orlib/crypto/hash"
)

type Address struct {
    IP net.IP
    Port uint16
    Nonce uint64
}

func (a *Address) Read(r io.Reader) (err error) {
    a.IP = make(net.IP, net.IPv6len)
    if err = butils.ReadFull(r, a.IP[:]); err != nil { return }
    if a.Port, err = butils.ReadUint16(r); err != nil { return }
    a.Nonce, err = butils.ReadUint64(r)
    return
}

func (a *Address) Write(w io.Writer) (err error) {
    a.IP = a.IP.To16()
    if a.IP == nil { return errors.New("Address ip not set") }
    if err = butils.WriteFull(w, a.IP.To16()[:]); err != nil { return }
    if err = butils.WriteUint16(w, a.Port); err != nil { return }
    return butils.WriteUint64(w, a.Nonce)
}

func (a *Address) String() string {
    return a.IP.String() + ":" + strconv.FormatUint(uint64(a.Port), 10) + " [nonce:" + strconv.FormatUint(a.Nonce, 16) + "]"
}

func (a *Address) IsInternal() bool {
    return a.Port == 0
}

func (a *Address) Id() *hash.ID {
    buf, err := butils.WriteToBytes(a)
    if err != nil { return nil }
    return hash.NewId(buf)
}
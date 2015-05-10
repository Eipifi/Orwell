package common
import (
    "io"
    "orwell/orlib/butils"
    "net"
    "strconv"
    "errors"
)

const NoPort Port = 0x00

type Port uint16

type Address struct {
    IP net.IP
    Port Port
}

func (a *Address) Read(r io.Reader) (err error) {
    a.IP = make(net.IP, net.IPv6len)
    if err = butils.ReadFull(r, a.IP[:]); err != nil { return }
    var port uint16
    if port, err = butils.ReadUint16(r); err != nil { return }
    a.Port = Port(port)
    return
}

func (a *Address) Write(w io.Writer) (err error) {
    a.IP = a.IP.To16()
    if a.IP == nil { return errors.New("Address ip not set") }
    if err = butils.WriteFull(w, a.IP.To16()[:]); err != nil { return }
    return butils.WriteUint16(w, uint16(a.Port))
}

func (a *Address) String() string {
    return a.IP.String() + ":" + strconv.FormatUint(uint64(a.Port), 10)
}
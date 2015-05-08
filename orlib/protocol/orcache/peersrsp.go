package orcache
import (
    "orwell/orlib/protocol/common"
    "io"
    "orwell/orlib/butils"
    "errors"
)

const MaxAddressesInPeersResp = 100
var ErrTooManyAddresses = errors.New("Too many addresses in PeersResp")

type PeersRsp struct {
    Addresses []common.Address
}

func (*PeersRsp) Code() byte { return 0x84 }

func (p *PeersRsp) Read(r io.Reader) (err error) {
    num, err := butils.ReadVarUint(r)
    if err != nil { return }
    if num > MaxAddressesInPeersResp { return ErrTooManyAddresses }
    for ; num > 0; num-- {
        a := common.Address{}
        if err = a.Read(r); err != nil { return }
        p.Addresses = append(p.Addresses, a)
    }
    return
}

func (p *PeersRsp) Write(w io.Writer) (err error) {
    if len(p.Addresses) > MaxAddressesInPeersResp { return ErrTooManyAddresses}
    if err = butils.WriteVarUint(w, uint64(len(p.Addresses))); err != nil { return }
    for _, a := range p.Addresses {
        if err = a.Write(w); err != nil { return }
    }
    return
}
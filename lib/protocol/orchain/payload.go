package orchain
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
    "errors"
)

const MAX_LABEL_LENGTH uint64 = 256

type Payload struct {
    label *[]byte
    ticket *foo.U256
    domain *Domain
    transfer *Transfer
}

func (p *Payload) Read(r io.Reader) (err error) {
    flag, err := butils.ReadByte(r)
    if err != nil { return }
    switch flag {
        case 0x00:
            return
        case 0x01:
            label, err := butils.ReadVarBytes(r, MAX_LABEL_LENGTH)
            p.label = &label
            return err
        case 0x02:
            p.ticket = &foo.U256{}
            return p.ticket.Read(r)
        case 0x03:
            p.domain = &Domain{}
            return p.domain.Read(r)
        case 0x04:
            p.transfer = &Transfer{}
            return p.transfer.Read(r)
        default:
            return errors.New("Unknown payload flag")
    }
}

func (p *Payload) Write(w io.Writer) (err error) {
    if p.label != nil {
        if err = butils.WriteByte(w, 0x01); err != nil { return }
        return butils.WriteVarBytes(w, *p.label, MAX_LABEL_LENGTH)
    }
    if p.ticket != nil {
        if err = butils.WriteByte(w, 0x02); err != nil { return }
        return p.ticket.Write(w)
    }
    if p.domain != nil {
        if err = butils.WriteByte(w, 0x03); err != nil { return }
        return p.domain.Write(w)
    }
    if p.transfer != nil {
        if err = butils.WriteByte(w, 0x04); err != nil { return }
        return p.transfer.Write(w)
    }
    return butils.WriteByte(w, 0x00)
}

func PayloadLabel2(label []byte) Payload {
    p := Payload{}
    p.label = &label
    return p
}

func PayloadLabelString(label string) Payload {
    return PayloadLabel2([]byte(label))
}


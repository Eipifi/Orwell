package orchain
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
    "errors"
)

const MAX_LABEL_LENGTH uint64 = 256

type Payload struct {
    Label *[]byte
    Ticket *foo.U256
    Domain *Domain
    Transfer *Transfer
}

func (p *Payload) Read(r io.Reader) (err error) {
    flag, err := butils.ReadByte(r)
    if err != nil { return }
    switch flag {
        case 0x00:
            return
        case 0x01:
            label, err := butils.ReadVarBytes(r, MAX_LABEL_LENGTH)
            p.Label = &label
            return err
        case 0x02:
            p.Ticket = &foo.U256{}
            return p.Ticket.Read(r)
        case 0x03:
            p.Domain = &Domain{}
            return p.Domain.Read(r)
        case 0x04:
            p.Transfer = &Transfer{}
            return p.Transfer.Read(r)
        default:
            return errors.New("Unknown payload flag")
    }
}

func (p *Payload) Write(w io.Writer) (err error) {
    if p.Label != nil {
        if err = butils.WriteByte(w, 0x01); err != nil { return }
        return butils.WriteVarBytes(w, *p.Label, MAX_LABEL_LENGTH)
    }
    if p.Ticket != nil {
        if err = butils.WriteByte(w, 0x02); err != nil { return }
        return p.Ticket.Write(w)
    }
    if p.Domain != nil {
        if err = butils.WriteByte(w, 0x03); err != nil { return }
        return p.Domain.Write(w)
    }
    if p.Transfer != nil {
        if err = butils.WriteByte(w, 0x04); err != nil { return }
        return p.Transfer.Write(w)
    }
    return butils.WriteByte(w, 0x00)
}

func PayloadLabel2(label []byte) Payload {
    return Payload{Label: &label}
}

func PayloadLabelString(label string) Payload {
    return PayloadLabel2([]byte(label))
}


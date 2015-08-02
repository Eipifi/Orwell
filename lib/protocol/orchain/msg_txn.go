package orchain
import (
    "io"
    "orwell/lib/butils"
)

const GET_TXNS_MAX_TXNS = 4096

type MsgGetTxns struct {}

func (*MsgGetTxns) Read(io.Reader) error { return nil }
func (m *MsgGetTxns) Write(io.Writer) error { return nil }

type MsgTxns struct {
    Transactions []Transaction
}

func (m *MsgTxns) Read(r io.Reader) (err error) {
    var num uint64
    if num, err = butils.ReadVarUint(r); err != nil { return }
    if num > GET_TXNS_MAX_TXNS { return ErrArrayTooLarge }
    m.Transactions = make([]Transaction, num)
    for i := 0; i < int(num); i += 1 {
        if err = m.Transactions[i].Read(r); err != nil { return }
    }
    return nil
}


func (m *MsgTxns) Write(w io.Writer) (err error) {
    num := uint64(len(m.Transactions))
    if num > GET_TXNS_MAX_TXNS { return ErrArrayTooLarge }
    if err = butils.WriteVarUint(w, num); err != nil { return }
    for i := 0; i < int(num); i += 1 {
        if err = m.Transactions[i].Write(w); err != nil { return }
    }
    return nil
}
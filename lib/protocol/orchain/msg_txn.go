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
    if err = butils.ReadSlice(r, GET_TXNS_MAX_TXNS, &m.Transactions); err != nil { return }
    return nil
}


func (m *MsgTxns) Write(w io.Writer) (err error) {
    if err = butils.WriteSlice(w, GET_TXNS_MAX_TXNS, m.Transactions); err != nil { return }
    return nil
}
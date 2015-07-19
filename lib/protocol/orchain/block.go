package orchain
import (
    "io"
    "orwell/lib/butils"
)

const BLOCK_TXN_MAX = 4096

type Block struct {
    Header Header
    Transactions []Transaction
}

func (b *Block) Read(r io.Reader) (err error) {
    if err = b.Header.Read(r); err != nil { return }

    var num uint64
    if num, err = butils.ReadVarUint(r); err != nil { return }
    if num > BLOCK_TXN_MAX { return ErrArrayTooLarge }
    b.Transactions = make([]Transaction, num)
    for i := 0; i < int(num); i += 1 {
        if err = b.Transactions[i].Read(r); err != nil { return }
    }

    return nil
}

func (b *Block) Write(w io.Writer) (err error) {
    if err = b.Header.Write(w); err != nil { return }

    num := uint64(len(b.Transactions))
    if num > BLOCK_TXN_MAX { return ErrArrayTooLarge }
    if err = butils.WriteVarUint(w, num); err != nil { return }
    for i := 0; i < int(num); i += 1 {
        if err = b.Transactions[i].Write(w); err != nil { return }
    }

    return nil
}
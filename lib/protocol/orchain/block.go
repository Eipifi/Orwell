package orchain
import (
    "io"
    "orwell/lib/butils"
)

const BLOCK_TXN_MAX = 4096
const BLOCK_DOMAIN_MAX = 4096 // TODO: think about it

type Block struct {
    Header Header
    Transactions []Transaction
    Domains []Domain
}

func (b *Block) Read(r io.Reader) (err error) {
    if err = b.Header.Read(r); err != nil { return }
    if err = butils.ReadSlice(r, BLOCK_TXN_MAX, &b.Transactions); err != nil { return }
    if err = butils.ReadSlice(r, BLOCK_DOMAIN_MAX, &b.Domains); err != nil { return }
    return nil
}

func (b *Block) Write(w io.Writer) (err error) {
    if err = b.Header.Write(w); err != nil { return }
    if err = butils.WriteSlice(w, BLOCK_TXN_MAX, b.Transactions); err != nil { return }
    if err = butils.WriteSlice(w, BLOCK_DOMAIN_MAX, b.Domains); err != nil { return }
    return nil
}
package orchain
import (
    "orwell/lib/crypto/merkle"
    "errors"
    "orwell/lib/butils"
)

type Block struct {
    Header Header
    Transactions []Transaction
}

func (b *Block) ComputeMerkleRoot() error {
    root, err := getMerkleRoot(b.Transactions)
    if err == nil { b.Header.MerkleRoot = root }
    return err
}

func getMerkleRoot(txns []Transaction) (hash butils.Uint256, err error) {
    n := len(txns)
    h := make([]butils.Uint256, n)
    for i := 0; i < n; i += 1 {
        h[i], err = txns[i].ID()
        if err != nil { return }
    }
    return merkle.Compute(h), nil
}

func (b *Block) CheckMerkleRoot() error {
    root, err := getMerkleRoot(b.Transactions)
    if err != nil { return err }
    if butils.Compare(b.Header.MerkleRoot, root) == 0 { return nil }
    return errors.New("Merkle roots do not match")
}
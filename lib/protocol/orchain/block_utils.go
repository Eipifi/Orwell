package orchain
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/merkle"
    "errors"
)

func (b *Block) ComputeMerkleRoot() error {
    root, err := getMerkleRoot(b.Transactions)
    if err == nil { b.Header.MerkleRoot = root }
    return err
}

func getMerkleRoot(txns []Transaction) (hash foo.U256, err error) {
    n := len(txns)
    h := make([]foo.U256, n)
    for i := 0; i < n; i += 1 {
        h[i], err = txns[i].TryID()
        if err != nil { return }
    }
    // TODO: also include domains
    return merkle.Compute(h), nil
}

func (b *Block) CheckMerkleRoot() error {
    root, err := getMerkleRoot(b.Transactions)
    if err != nil { return err }
    if foo.Compare(b.Header.MerkleRoot, root) == 0 { return nil }
    return errors.New("Merkle roots do not match")
}

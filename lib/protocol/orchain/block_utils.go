package orchain
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/merkle"
    "errors"
)

func (b *Block) ComputeMerkleRoot() error {
    root, err := getMerkleRoot(b.Transactions, b.Domains)
    if err == nil { b.Header.MerkleRoot = root }
    return err
}

func getMerkleRoot(txns []Transaction, domains []Domain) (hash foo.U256, err error) {
    var ids []foo.U256
    for _, txn := range txns {
        ids = append(ids, txn.ID())
    }
    for _, dmn := range domains {
        ids = append(ids, dmn.ID())
    }
    return merkle.Compute(ids), nil
}

func (b *Block) CheckMerkleRoot() error {
    root, err := getMerkleRoot(b.Transactions, b.Domains)
    if err != nil { return err }
    if foo.Compare(b.Header.MerkleRoot, root) == 0 { return nil }
    return errors.New("Merkle roots do not match")
}

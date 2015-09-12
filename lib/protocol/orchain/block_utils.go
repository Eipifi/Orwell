package orchain
import (
    "orwell/lib/foo"
    "orwell/lib/crypto/merkle"
    "errors"
)

func (b *Block) ComputeMerkleRoot() {
    b.Header.MerkleRoot = getMerkleRoot(b.Transactions, b.Domains)
}

func getMerkleRoot(txns []Transaction, domains []Domain) (hash foo.U256) {
    var ids []foo.U256
    for _, txn := range txns {
        ids = append(ids, txn.ID())
    }
    for _, dmn := range domains {
        ids = append(ids, dmn.ID())
    }
    return merkle.Compute(ids)
}

func (b *Block) CheckMerkleRoot() error {
    root := getMerkleRoot(b.Transactions, b.Domains)
    if foo.Compare(b.Header.MerkleRoot, root) == 0 { return nil }
    return errors.New("Merkle roots do not match")
}

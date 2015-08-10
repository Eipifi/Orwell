package db
import "orwell/lib/protocol/orchain"

func CheckBlockMerkleRoot(t *Tx, b *orchain.Block) (err error) {
    return b.CheckMerkleRoot()
}



package blockstore
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/butils"
)

type BlockStorage interface {
    Push(*orchain.Block) error
    Pop()
    Length() uint64
    Head() butils.Uint256
}

// This database is responsible for keeping the storage in consistent state, with atomic updates.
// The higher layer must check if the block is legal.
type Database interface {

    // Basic info methods
    Length() uint64
    Head() butils.Uint256

    // Write operations
    PutBlock(*orchain.Block) error
    PopBlock() error

    // Read operations
    GetHeaderByID(butils.Uint256) *orchain.Header
    GetHeaderByNum(uint64) *orchain.Header
    GetTransaction(butils.Uint256) *orchain.Transaction
    GetTransactions(butils.Uint256) []butils.Uint256
    GetBill(orchain.BillNumber) *orchain.Bill

}
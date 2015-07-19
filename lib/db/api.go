package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
)

type DB interface {
    Push(*orchain.Block) error
    Pop()
    State() *State
    GetBlockByID(foo.U256) *orchain.Block
    GetHeaderByID(foo.U256) *orchain.Header
    GetHeaderByNum(uint64) *orchain.Header
    GetNumByID(foo.U256) *uint64
    GetIDByNum(uint64) *foo.U256
}

// Storage is responsible for persisting the data in a consistent state, with atomic updates.
// The higher layer must check if the block is legal.
type Storage interface {

    // Basic info methods
    State() *State

    // Write operations
    PutBlock(*orchain.Block) error
    PopBlock() error

    // Read operations
    GetHeaderByID(foo.U256) *orchain.Header
    GetHeaderByNum(uint64) *orchain.Header
    GetIDByNum(uint64) *foo.U256
    GetNumByID(foo.U256) *uint64
    GetTransaction(foo.U256) *orchain.Transaction
    GetTransactions(foo.U256) []foo.U256
    GetBill(orchain.BillNumber) *orchain.Bill
}
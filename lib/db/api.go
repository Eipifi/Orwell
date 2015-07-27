package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
)

type DB interface {
    Push(*orchain.Block) error
    Pop()
    State() *State
    Difficulty() foo.U256

    GetBlockByID(foo.U256) *orchain.Block
    GetHeaderByNum(uint64) *orchain.Header

    GetIDByNum(uint64) *foo.U256
    GetNumByID(foo.U256) *uint64

    GetBills(foo.U256) []orchain.BillNumber
    GetBill(orchain.BillNumber) *orchain.Bill
}
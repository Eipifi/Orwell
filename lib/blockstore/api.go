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

type Database interface {
    StoreHead(butils.Uint256, uint64)
    FetchHead() (butils.Uint256, uint64)

    StoreHeader(*orchain.Header) error
    FetchHeader(butils.Uint256) *orchain.Header
    RemoveHeader(butils.Uint256)

    StoreBlockTransactionIDs(butils.Uint256, []butils.Uint256)
    FetchBlockTransactionIDs(butils.Uint256) []butils.Uint256
    RemoveBlockTransactionIDs(butils.Uint256)

    StoreTransaction(*orchain.Transaction) error
    FetchTransaction(butils.Uint256) *orchain.Transaction
    RemoveTransaction(butils.Uint256)

    StoreUnspentBill(orchain.BillNumber, orchain.Bill)
    FetchUnspentBill(orchain.BillNumber) *orchain.Bill
    SpendBill(orchain.BillNumber)
}
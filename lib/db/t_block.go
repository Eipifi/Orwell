package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/butils"
)

var BUCKET_HEADER     = []byte("header")
var BUCKET_HID_NUM    = []byte("hid_num")
var BUCKET_NUM_HID    = []byte("num_hid")

func (t *Tx) PutBlock(block *orchain.Block, num uint64) {
    bid := block.Header.ID()
    t.Put(BUCKET_NUM_HID, butils.Uint64ToBytes(num), bid[:])
    t.Put(BUCKET_HID_NUM, bid[:], butils.Uint64ToBytes(num))
    t.Write(BUCKET_HEADER, bid[:], &block.Header)
    t.PutTransactionsFromBlock(bid, block.Transactions)
    t.PutRegisteredDomainsFromBlock(bid, block.Domains)
}

func (t *Tx) GetBlock(id foo.U256) (b *orchain.Block) {
    h := t.GetHeaderByID(id)
    if h == nil { return }
    b = &orchain.Block{}
    b.Header = *h
    b.Transactions = t.GetTransactionsFromBlock(id)
    b.Domains = t.GetDomainsFromBlock(id)
    return
}

func (t *Tx) GetHeaderByID(id foo.U256) (h *orchain.Header) {
    h = &orchain.Header{}
    if t.Read(BUCKET_HEADER, id[:], h) { return }
    return nil
}

func (t *Tx) GetHeaderByNum(num uint64) (h *orchain.Header) {
    hid := t.Get(BUCKET_NUM_HID, butils.Uint64ToBytes(num))
    if hid == nil { return nil }
    h = &orchain.Header{}
    if t.Read(BUCKET_HEADER, hid, h) { return }
    return nil
}

func (t *Tx) GetIDByNum(num uint64) (id *foo.U256) {
    id = &foo.U256{}
    if t.Read(BUCKET_NUM_HID, butils.Uint64ToBytes(num), id) { return }
    return nil
}

func (t *Tx) GetNumByID(id foo.U256) (num *uint64) {
    data := t.Get(BUCKET_HID_NUM, id[:])
    if data == nil { return nil }
    res := butils.BytesToUint64(data)
    return &res
}
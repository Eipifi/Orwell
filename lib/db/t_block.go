package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/butils"
    "bytes"
    "io"
    "orwell/lib/utils"
    "github.com/boltdb/bolt"
)

var BUCKET_HEADER     = []byte("header")
var BUCKET_HID_NUM    = []byte("hid_num")
var BUCKET_NUM_HID    = []byte("num_hid")
var BUCKET_TXN_LIST   = []byte("txn_list")

func PutBlock(t *bolt.Tx, block *orchain.Block, num uint64) {
    bid := block.Header.ID()
    Put(t, BUCKET_NUM_HID, butils.Uint64ToBytes(num), bid[:])
    Put(t, BUCKET_HID_NUM, bid[:], butils.Uint64ToBytes(num))
    Write(t, BUCKET_HEADER, bid[:], &block.Header)
    buf := &bytes.Buffer{}
    for _, txn := range block.Transactions {
        tid := txn.ID()
        tid.Write(buf)
        PutTransaction(t, &txn)
    }
    Put(t, BUCKET_TXN_LIST, bid[:], buf.Bytes())
}

func GetBlock(t *bolt.Tx, id foo.U256) (b *orchain.Block) {
    h := GetHeaderByID(t, id)
    if h == nil { return }
    b = &orchain.Block{}
    b.Header = *h
    tids := Get(t, BUCKET_TXN_LIST, id[:])
    buf := bytes.NewBuffer(tids)
    for {
        var tid foo.U256
        err := tid.Read(buf)
        if err == io.EOF { break }
        utils.Ensure(err)
        txn := GetTransaction(t, tid)
        b.Transactions = append(b.Transactions, *txn)
    }
    return
}

func GetHeaderByID(t *bolt.Tx, id foo.U256) (h *orchain.Header) {
    h = &orchain.Header{}
    if Read(t, BUCKET_HEADER, id[:], h) { return }
    return nil
}

func GetHeaderByNum(t *bolt.Tx, num uint64) (h *orchain.Header) {
    hid := Get(t, BUCKET_NUM_HID, butils.Uint64ToBytes(num))
    if hid == nil { return nil }
    h = &orchain.Header{}
    if Read(t, BUCKET_HEADER, hid, h) { return }
    return nil
}

func GetIDByNum(t *bolt.Tx, num uint64) (id *foo.U256) {
    id = &foo.U256{}
    if Read(t, BUCKET_NUM_HID, butils.Uint64ToBytes(num), id) { return }
    return nil
}

func GetNumByID(t *bolt.Tx, id foo.U256) (num *uint64) {
    data := Get(t, BUCKET_HID_NUM, id[:])
    if data == nil { return nil }
    res := butils.BytesToUint64(data)
    return &res
}

func DelBlock(t *bolt.Tx, id foo.U256) {
    num_bytes := Get(t, BUCKET_HID_NUM, id[:])
    tids := Get(t, BUCKET_TXN_LIST, id[:])
    Del(t, BUCKET_NUM_HID, num_bytes)
    Del(t, BUCKET_HID_NUM, id[:])
    Del(t, BUCKET_HEADER, id[:])
    Del(t, BUCKET_TXN_LIST, id[:])
    buf := bytes.NewBuffer(tids)
    for {
        var tid foo.U256
        err := tid.Read(buf)
        if err == io.EOF { break }
        utils.Ensure(err)
        DelTransaction(t, tid)
    }
    return
}
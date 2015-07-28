package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/butils"
    "orwell/lib/utils"
    "bytes"
)

type BillStatus byte

const (
    NONEXISTENT BillStatus = 0  // no transaction object in storage
    UNSPENT     BillStatus = 1  // transaction exists, bill not spent (flag is set)
    SPENT       BillStatus = 2  // transaction exists, bill spent (no flag)
)

var BUCKET_UNSPENT = []byte("unspent")  // NUM -> 0xFF
var BUCKET_OWNED = []byte("owned")      // OWNER | NUM -> 0xFF
var FLAG = []byte{0xFF}

func (t *Tx) GetBill(num *orchain.BillNumber) *orchain.Bill {
    txn := t.GetTransaction(num.Txn)
    if txn == nil { return nil }
    if uint64(len(txn.Outputs)) <= num.Index { return nil }
    return &txn.Outputs[num.Index]
}

func (t *Tx) GetBillStatus(num *orchain.BillNumber) BillStatus {
    bill := t.GetBill(num)
    if bill == nil { return NONEXISTENT }
    if t.Get(BUCKET_UNSPENT, butils.ToBytes(num)) == nil { return SPENT }
    return UNSPENT
}

func (t *Tx) SetBillStatus(num *orchain.BillNumber, status BillStatus) {
    bill := t.GetBill(num)
    key := butils.ToBytes(num)
    if status == UNSPENT {
        utils.Assert(bill != nil)
        t.Put(BUCKET_UNSPENT, key, FLAG)
        t.Put(BUCKET_OWNED, utils.Cat(bill.Target[:], key), FLAG)
    } else {
        if bill == nil {
            if status == NONEXISTENT { return }
            panic("Tried to spend a nonexistent bill")
        }
        t.Del(BUCKET_UNSPENT, key)
        t.Del(BUCKET_OWNED, utils.Cat(bill.Target[:], key))
    }
}

func (t *Tx) GetUnspentBillsByWallet(wallet foo.U256) (res []orchain.BillNumber) {
    c := t.tx.Bucket(BUCKET_OWNED).Cursor()
    prefix := wallet[:]
    for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
        num := orchain.BillNumber{}
        butils.ReadAllInto(&num, k[foo.U256_BYTES:])
        res = append(res, num)
    }
    return
}

func (t *Tx) GetBalance(wallet foo.U256) (sum uint64) {
    for _, inp := range t.GetUnspentBillsByWallet(wallet) {
        sum += t.GetBill(&inp).Value
    }
    return
}
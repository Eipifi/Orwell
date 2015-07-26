package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/butils"
    "orwell/lib/utils"
    "bytes"
    "github.com/boltdb/bolt"
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

func GetBill(t *bolt.Tx, num *orchain.BillNumber) *orchain.Bill {
    txn := GetTransaction(t, num.Txn)
    if txn == nil { return nil }
    if uint64(len(txn.Outputs)) <= num.Index { return nil }
    return &txn.Outputs[num.Index]
}

func GetBillStatus(t *bolt.Tx, num *orchain.BillNumber) BillStatus {
    bill := GetBill(t, num)
    if bill == nil { return NONEXISTENT }
    if Get(t, BUCKET_UNSPENT, butils.ToBytes(num)) == nil { return SPENT }
    return UNSPENT
}

func SetBillStatus(t *bolt.Tx, num *orchain.BillNumber, status BillStatus) {
    bill := GetBill(t, num)
    key := butils.ToBytes(num)
    if status == UNSPENT {
        utils.Assert(bill != nil)
        Put(t, BUCKET_UNSPENT, key, FLAG)
        Put(t, BUCKET_OWNED, utils.Cat(bill.Target[:], key), FLAG)
    } else {
        if bill == nil {
            if status == NONEXISTENT { return }
            panic("Tried to spend a nonexistent bill")
        }
        Del(t, BUCKET_UNSPENT, key)
        Del(t, BUCKET_OWNED, utils.Cat(bill.Target[:], key))
    }
}

func GetUnspentBillsByUser(t *bolt.Tx, user foo.U256) (res []orchain.BillNumber) {
    c := t.Bucket(BUCKET_OWNED).Cursor()
    prefix := user[:]
    for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
        num := orchain.BillNumber{}
        butils.ReadAllInto(&num, k[foo.U256_BYTES:])
        res = append(res, num)
    }
    return
}
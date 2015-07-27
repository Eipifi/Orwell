package db
import (
    "github.com/boltdb/bolt"
    "orwell/lib/utils"
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
)

type BoltDB struct {
    db *bolt.DB
}

func open(path string) (r *BoltDB) {
    db, err := bolt.Open(path, 0600, nil)
    utils.Ensure(err)
    r = &BoltDB{db: db}
    utils.Ensure(r.db.Update(initialize))
    return
}

func initialize(tx *bolt.Tx) (err error) {
    if _, err = tx.CreateBucketIfNotExists(BUCKET_INFO); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_TXN); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_UNSPENT); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_OWNED); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_HEADER); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_HID_NUM); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_NUM_HID); err != nil { return }
    if _, err = tx.CreateBucketIfNotExists(BUCKET_TXN_LIST); err != nil { return }
    if GetState(tx) == nil {
        PutState(tx, &State{})
        utils.Ensure(PushBlock(tx, GenesisBlock()))
    }
    return
}

func (b *BoltDB) Push(block *orchain.Block) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        return PushBlock(tx, block)
    })
}

func (b *BoltDB) Pop() {
    b.db.Update(func(tx *bolt.Tx) error {
        PopBlock(tx)
        return nil
    })
}

func (b *BoltDB) State() (state *State) {
    b.db.View(func(tx *bolt.Tx) error {
        state = GetState(tx)
        return nil
    })
    return
}

func (b *BoltDB) Difficulty() (difficulty foo.U256) {
    b.db.View(func(tx *bolt.Tx) error {
        difficulty = Difficulty(tx)
        return nil
    })
    return
}

func (b *BoltDB) GetBlockByID(id foo.U256) (block *orchain.Block) {
    b.db.View(func(tx *bolt.Tx) error {
        block = GetBlock(tx, id)
        return nil
    })
    return
}

func (b *BoltDB) GetHeaderByNum(num uint64) (header *orchain.Header) {
    b.db.View(func(tx *bolt.Tx) error {
        header = GetHeaderByNum(tx, num)
        return nil
    })
    return
}

func (b *BoltDB) GetIDByNum(num uint64) (id *foo.U256) {
    b.db.View(func(tx *bolt.Tx) error {
        id = GetIDByNum(tx, num)
        return nil
    })
    return
}

func (b *BoltDB) GetNumByID(id foo.U256) (num *uint64) {
    b.db.View(func(tx *bolt.Tx) error {
        num = GetNumByID(tx, id)
        return nil
    })
    return
}

func (b *BoltDB) GetBills(wallet foo.U256) (res []orchain.BillNumber) {
    b.db.View(func(tx *bolt.Tx) error {
        res = GetUnspentBillsByUser(tx, wallet)
        return nil
    })
    return
}

func (b *BoltDB) GetBill(address orchain.BillNumber) (res *orchain.Bill) {
    b.db.View(func(tx *bolt.Tx) error {
        res = GetBill(tx, &address)
        return nil
    })
    return
}


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var instance DB

func GetDB() DB { return instance }

func Initialize(path string) {
    instance = open(path)
}



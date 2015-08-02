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

type Tx struct {
    tx *bolt.Tx
}

func open(path string) (r *BoltDB) {
    db, err := bolt.Open(path, 0600, nil)
    utils.Ensure(err)
    r = &BoltDB{db: db}
    utils.Ensure(r.UpdateE(initialize))
    return
}

func initialize(t *Tx) (err error) {
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_INFO); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_TXN); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_UNSPENT); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_OWNED); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_HEADER); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_HID_NUM); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_NUM_HID); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_TXN_LIST); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_TXN_UNCONFIRMED); err != nil { return }
    if t.GetState() == nil {
        t.PutState(&State{})
        utils.Ensure(t.PushBlock(orchain.GenesisBlock()))
    }
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BoltDB) Push(block *orchain.Block) error {
    return b.UpdateE(func(t *Tx) error {
        return t.PushBlock(block)
    })
}

func (b *BoltDB) Pop() {
    b.Update(func(t *Tx) {
        t.PopBlock()
    })
}

func (b *BoltDB) State() (state *State) {
    b.View(func(t *Tx) {
        state = t.GetState()
    })
    return
}

func (b *BoltDB) Difficulty() (difficulty foo.U256) {
    b.View(func(t *Tx) {
        difficulty = t.GetDifficulty()
    })
    return
}

func (b *BoltDB) GetBlockByID(id foo.U256) (block *orchain.Block) {
    b.View(func(t *Tx) {
        block = t.GetBlock(id)
    })
    return
}

func (b *BoltDB) GetHeaderByNum(num uint64) (header *orchain.Header) {
    b.View(func(t *Tx) {
        header = t.GetHeaderByNum(num)
    })
    return
}

func (b *BoltDB) GetIDByNum(num uint64) (id *foo.U256) {
    b.View(func(t *Tx) {
        id = t.GetIDByNum(num)
    })
    return
}

func (b *BoltDB) GetNumByID(id foo.U256) (num *uint64) {
    b.View(func(t *Tx) {
        num = t.GetNumByID(id)
    })
    return
}

func (b *BoltDB) GetBills(wallet foo.U256) (res []orchain.BillNumber) {
    b.View(func(t *Tx) {
        res = t.GetUnspentBillsByWallet(wallet)
    })
    return
}

func (b *BoltDB) GetBill(address orchain.BillNumber) (res *orchain.Bill) {
    b.View(func(t *Tx) {
        res = t.GetBill(&address)
    })
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BoltDB) UpdateE(f func(*Tx) error) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        return f(&Tx{tx})
    })
}

func (b *BoltDB) ViewE(f func(*Tx) error) error {
    return b.db.View(func(tx *bolt.Tx) error {
        return f(&Tx{tx})
    })
}

func (b *BoltDB) Update(f func(*Tx)) {
    utils.Ensure(b.db.Update(func(tx *bolt.Tx) error {
        f(&Tx{tx})
        return nil
    }))
}

func (b *BoltDB) View(f func(*Tx)) {
    utils.Ensure(b.db.View(func(tx *bolt.Tx) error {
        f(&Tx{tx})
        return nil
    }))
}


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var instance *BoltDB

func Get() *BoltDB { return instance }

func Initialize(path string) {
    instance = open(path)
}



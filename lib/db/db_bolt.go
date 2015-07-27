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
    utils.Ensure(r.Update(initialize))
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
    if t.GetState() == nil {
        t.PutState(&State{})
        utils.Ensure(t.PushBlock(GenesisBlock()))
    }
    return
}

func (b *BoltDB) Push(block *orchain.Block) error {
    return b.Update(func(t *Tx) error {
        return t.PushBlock(block)
    })
}

func (b *BoltDB) Pop() {
    b.Update(func(t *Tx) error {
        t.PopBlock()
        return nil
    })
}

func (b *BoltDB) State() (state *State) {
    b.View(func(t *Tx) error {
        state = t.GetState()
        return nil
    })
    return
}

func (b *BoltDB) Difficulty() (difficulty foo.U256) {
    b.View(func(t *Tx) error {
        difficulty = t.Difficulty()
        return nil
    })
    return
}

func (b *BoltDB) GetBlockByID(id foo.U256) (block *orchain.Block) {
    b.View(func(t *Tx) error {
        block = t.GetBlock(id)
        return nil
    })
    return
}

func (b *BoltDB) GetHeaderByNum(num uint64) (header *orchain.Header) {
    b.View(func(t *Tx) error {
        header = t.GetHeaderByNum(num)
        return nil
    })
    return
}

func (b *BoltDB) GetIDByNum(num uint64) (id *foo.U256) {
    b.View(func(t *Tx) error {
        id = t.GetIDByNum(num)
        return nil
    })
    return
}

func (b *BoltDB) GetNumByID(id foo.U256) (num *uint64) {
    b.View(func(t *Tx) error {
        num = t.GetNumByID(id)
        return nil
    })
    return
}

func (b *BoltDB) GetBills(wallet foo.U256) (res []orchain.BillNumber) {
    b.View(func(t *Tx) error {
        res = t.GetUnspentBillsByUser(wallet)
        return nil
    })
    return
}

func (b *BoltDB) GetBill(address orchain.BillNumber) (res *orchain.Bill) {
    b.View(func(t *Tx) error {
        res = t.GetBill(&address)
        return nil
    })
    return
}

func (b *BoltDB) Update(f func(*Tx) error) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        return f(&Tx{tx})
    })
}

func (b *BoltDB) View(f func(*Tx) error) error {
    return b.db.View(func(tx *bolt.Tx) error {
        return f(&Tx{tx})
    })
}


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var instance DB

func Get() DB { return instance }

func Initialize(path string) {
    instance = open(path)
}



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
    session *foo.U256
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
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_ROLLBACK); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_TXN); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_UNSPENT); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_OWNED); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_HEADER); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_NUM_HID); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_HID_NUM); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_TXN_LIST); err != nil { return }
    if _, err = t.tx.CreateBucketIfNotExists(BUCKET_TXN_UNCONFIRMED); err != nil { return }
    state := t.GetState()
    if state.Length == 0 {
        utils.Ensure(t.PushBlock(orchain.GenesisBlock()))
    }
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BoltDB) UpdateE(f func(*Tx) error) error {
    return b.db.Update(func(tx *bolt.Tx) error {
        return f(&Tx{tx, nil})
    })
}

func (b *BoltDB) ViewE(f func(*Tx) error) error {
    return b.db.View(func(tx *bolt.Tx) error {
        return f(&Tx{tx, nil})
    })
}

func (b *BoltDB) Update(f func(*Tx)) {
    utils.Ensure(b.db.Update(func(tx *bolt.Tx) error {
        f(&Tx{tx, nil})
        return nil
    }))
}

func (b *BoltDB) View(f func(*Tx)) {
    utils.Ensure(b.db.View(func(tx *bolt.Tx) error {
        f(&Tx{tx, nil})
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



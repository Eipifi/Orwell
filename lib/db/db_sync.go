package db
import (
    "sync"
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
)

type SyncDB struct {
    mtx sync.RWMutex
    db DB
}

func NewSyncDB(d DB) *SyncDB {
    r := &SyncDB{}
    r.db = d
    return r
}

func (d *SyncDB) State() *State {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    return d.db.State()
}

func (d *SyncDB) GetBlockByID(id foo.U256) *orchain.Block {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    return d.db.GetBlockByID(id)
}

func (d *SyncDB) GetHeaderByID(id foo.U256) *orchain.Header {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    return d.db.GetHeaderByID(id)
}

func (d *SyncDB) GetHeaderByNum(num uint64) *orchain.Header {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    return d.db.GetHeaderByNum(num)
}

func (d *SyncDB) GetNumByID(id foo.U256) *uint64 {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    return d.db.GetNumByID(id)
}

func (d *SyncDB) GetIDByNum(num uint64) *foo.U256 {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    return d.db.GetIDByNum(num)
}

func (d *SyncDB) Push(b *orchain.Block) (err error) {
    d.mtx.Lock()
    defer d.mtx.Unlock()
    return d.db.Push(b)
}

func (d *SyncDB) Pop() {
    d.mtx.Lock()
    defer d.mtx.Unlock()
    d.db.Pop()
}
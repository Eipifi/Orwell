package db
import (
    "orwell/lib/butils"
    "orwell/lib/utils"
)

func (t *Tx) Put(bucket, key, value []byte) {
    utils.Ensure(t.tx.Bucket(bucket).Put(key, value))
}

func (t *Tx) Get(bucket, key []byte) []byte {
    return t.tx.Bucket(bucket).Get(key)
}

func (t *Tx) Del(bucket, key []byte) {
    utils.Ensure(t.tx.Bucket(bucket).Delete(key))
}

func (t *Tx) Read(bucket, key []byte, target butils.Readable) bool {
    data := t.Get(bucket, key)
    if data == nil { return false }
    utils.Ensure(butils.ReadAllInto(target, data))
    return true
}

func (t *Tx) Write(bucket, key []byte, target butils.Writable) {
    t.Put(bucket, key, butils.ToBytes(target))
}
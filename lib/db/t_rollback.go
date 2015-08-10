package db
import (
    "orwell/lib/butils"
    "orwell/lib/foo"
    "bytes"
    "orwell/lib/utils"
)

var BUCKET_ROLLBACK = []byte("rollback")

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Tx) Put(bucket, key, value []byte) {
    t.preserve(bucket, key)
    t.RawPut(bucket, key, value)
}

func (t *Tx) Del(bucket, key []byte) {
    t.preserve(bucket, key)
    t.RawDel(bucket, key)
}

func (t *Tx) Get(bucket, key []byte) []byte {
    return t.tx.Bucket(bucket).Get(key)
}

func (t *Tx) preserve(bucket, key []byte) {
    rkey := utils.Cat(t.session[:], bucket, key)
    if t.Get(BUCKET_ROLLBACK, rkey) == nil {
        // no snapshot of the field was made, let's make one now
        prev := &RollbackOp{Bucket: bucket, Key: key, Value: t.Get(bucket, key)}
        t.RawWrite(BUCKET_ROLLBACK, rkey, prev)
    }
}

func (t *Tx) Rollback(session foo.U256) {
    c := t.tx.Bucket([]byte(BUCKET_ROLLBACK)).Cursor()

    prefix := session[:]
    for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
        utils.Ensure(c.Delete())
        prev := &RollbackOp{}
        utils.Ensure(butils.ReadAllInto(prev, v))
        if prev.Value == nil {
            t.RawDel(prev.Bucket, prev.Key)
        } else {
            t.RawPut(prev.Bucket, prev.Key, prev.Value)
        }
    }
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Tx) Read(bucket, key []byte, target butils.Readable) bool {
    data := t.Get(bucket, key)
    if data == nil { return false }
    utils.Ensure(butils.ReadAllInto(target, data))
    return true
}

func (t *Tx) Write(bucket, key []byte, target butils.Writable) {
    t.Put(bucket, key, butils.ToBytes(target))
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


func (t *Tx) RawPut(bucket, key, value []byte) {
    utils.Ensure(t.tx.Bucket(bucket).Put(key, value))
}

func (t *Tx) RawDel(bucket, key []byte) {
    utils.Ensure(t.tx.Bucket(bucket).Delete(key))
}

func (t *Tx) RawWrite(bucket, key []byte, target butils.Writable) {
    t.RawPut(bucket, key, butils.ToBytes(target))
}
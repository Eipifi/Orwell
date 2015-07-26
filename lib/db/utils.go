package db
import (
    "orwell/lib/butils"
    "orwell/lib/utils"
    "github.com/boltdb/bolt"
)

func Put(t *bolt.Tx, bucket, key, value []byte) {
    utils.Ensure(t.Bucket(bucket).Put(key, value))
}

func Get(t *bolt.Tx, bucket, key []byte) []byte {
    return t.Bucket(bucket).Get(key)
}

func Del(t *bolt.Tx, bucket, key []byte) {
    utils.Ensure(t.Bucket(bucket).Delete(key))
}

func Read(t *bolt.Tx, bucket, key []byte, target butils.Readable) bool {
    data := Get(t, bucket, key)
    if data == nil { return false }
    utils.Ensure(butils.ReadAllInto(target, data))
    return true
}

func Write(t *bolt.Tx, bucket, key []byte, target butils.Writable) {
    Put(t, bucket, key, butils.ToBytes(target))
}
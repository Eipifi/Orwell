package blockstore

var instance BlockStorage

func Initialize(path string) {
    db := Open(path)
    instance = NewBlockStore(db)
}

func Get() BlockStorage {
    return instance
}
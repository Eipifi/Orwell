package db

var instance DB

func Initialize(path string) {
    instance = NewSyncDB(NewDB(OpenLDBStorage(path)))
}

func Get() DB {
    // TODO: return a mutex-safe wrapped instance
    return instance
}
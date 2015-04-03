package main

type Env struct {
    Cache Cache
    Locker TokenLocker
    Manager ConnectionManager
}

func main() {
    serve(1984, RunPeer)
}
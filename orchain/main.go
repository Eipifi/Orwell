package main
import (
    "net"
    "code.google.com/p/leveldb-go/leveldb"
    "fmt"
    "github.com/mitchellh/go-homedir"
)

func main2() {
    listener, err := net.Listen("tcp", ":1984")
    if err != nil { return }
    for {
        c, err := listener.Accept()
        if err != nil { return }
        p := NewPeer(c)
        go p.Lifecycle()
    }
}

func main() {
    path, err := homedir.Expand("~/.orchain/db/")
    db, err := leveldb.Open(path, nil)
    fmt.Println(err)
    err = db.Set([]byte("qwe"), []byte("rty"), nil)
    fmt.Println(err)
    db.Close()
}
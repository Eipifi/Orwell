package main
import (
    "fmt"
    "net"
)

func main() {
    if err := run(); err != nil {
        fmt.Println(err)
    }
}

func run() error {
    mgr := NewManagerImpl()
    socket, err := net.Listen("tcp", ":1984")
    if err != nil { return err }

    for {
        conn, err := socket.Accept()
        if err != nil { return err }
        NewPeer(conn, mgr)
    }
}

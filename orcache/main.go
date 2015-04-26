package main
import (
    "fmt"
    "net"
)

func main() {
    err := run()
    if err != nil {
        fmt.Println(err)
    }
}

func run() error {
    socket, err := net.Listen("tcp", ":1984")
    if err != nil { return err }

    for {
        conn, err := socket.Accept()
        if err != nil { return err }
        HandleConnection(conn)
    }
}

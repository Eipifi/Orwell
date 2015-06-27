package main
import (
    "net"
    "fmt"
)

func runServer(port uint16) {
    fmt.Printf("Connecting to port %v\n", port)
    listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
    if err != nil { return }
    for {
        c, err := listener.Accept()
        if err != nil { return }
        go Connect(c)
    }
}
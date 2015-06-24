package main
import "net"


func main() {
    listener, err := net.Listen("tcp", ":1984")
    if err != nil { return }
    for {
        c, err := listener.Accept()
        if err != nil { return }
        go HandleConnection(c)
    }
}
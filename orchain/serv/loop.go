package serv
import (
    "net"
    "fmt"
)

func RunServer(port int) {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
    if err != nil { return }
    for {
        c, err := listener.Accept()
        if err != nil { return }
        go Talk(c)
    }
}
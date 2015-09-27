package obp
import (
    "net"
    "fmt"
)


func Serve(port int, handler func(net.Conn)) error {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
    if err != nil { return err }
    for {
        c, err := listener.Accept()
        if err != nil { return err }
        go handler(c)
    }
}
package serv
import (
    "net"
    "fmt"
    "orwell/lib/utils"
)

func RunServer(port int) {
    utils.Ensure(run(port))
}

func run(port int) error  {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
    if err != nil { return err }
    for {
        c, err := listener.Accept()
        if err != nil { return err }
        go Talk(c)
    }
}
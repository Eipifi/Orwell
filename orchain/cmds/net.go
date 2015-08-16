package cmds
import (
    "orwell/lib/fcli"
    "fmt"
    "net"
    "orwell/orchain/serv"
)


func NetAddHandler(address string) fcli.Result {
    fmt.Printf("Connecting to %v \n", address)
    conn, err := net.Dial("tcp", address)
    if err == nil {
        go serv.Talk(conn)
    }
    return err
}
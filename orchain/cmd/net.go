package cmd
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

func NetStatsHandler() fcli.Result {
    fmt.Println("Connected peers:")
    for _, peer := range serv.ConnMgr().GetAllPeers() {
        fmt.Println(peer.Info())
    }
    return nil
}
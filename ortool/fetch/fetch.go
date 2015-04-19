package fetch
import (
    "net"
    "fmt"
    "flag"
    "orwell/orlib/card"
    "orwell/orlib/protocol/types"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/comm"
)

const Usage = `usage: ortool fetch [--from ip:port] <address>

Asks the server for a card.

`
func Main(args []string) {
    if len(args) == 0 {
        fmt.Println("ID not provided.")
        return
    }

    var id *types.ID
    var err error
    if id, err = types.HexToID(args[len(args)-1]); err != nil {
        fmt.Println("Error: ", err)
        return
    }

    args = args[:len(args)-1]

    fs := flag.NewFlagSet("fetch", flag.ExitOnError)
    fSr := fs.String("from", "127.0.0.1:1984", "Orcache server address")
    fs.Parse(args)

    var conn net.Conn
    if conn, err = net.Dial("tcp", *fSr); err != nil {
        fmt.Println("Error:", err)
        return
    }

    var ms *orcache.OrcacheMessenger
    if ms, err = orcache.NewOrcacheMessenger(conn, "Ortool", nil); err != nil {
        fmt.Println("Error:", err)
        return
    }

    if err = ms.Write(&orcache.Get{types.RandomToken(), types.MaxTTLValue, id, 0}); err != nil {
        fmt.Println("Error:", err)
        return
    }

    var msg comm.Msg
    if msg, err = ms.ReadAny(); err != nil {
        fmt.Println("Error:", err)
        return
    }

    switch m := msg.(type) {
        case *orcache.CardFound:
            if c, e := card.Unmarshal(m.Card); e != nil {
                s := string(c.Payload.MarshalJSON())
                fmt.Println(s)
            } else {
                fmt.Println("Invalid card format or signature received.")
                return
            }
        case *orcache.CardNotFound:
            fmt.Printf("Failed to fetch the card. TTL=%d\n", m.TTL)
    }
}
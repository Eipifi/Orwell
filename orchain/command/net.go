package command
import (
    "fmt"
    "errors"
    "net"
    "orwell/orchain/serv"
)

type NetCmd struct{}

func (*NetCmd) Name() string {
    return "net"
}

func (n *NetCmd) Run(args []string) error {
    if len(args) == 0 {
        return n.doInfo()
    }
    c := args[0]
    r := args[1:]

    if c == "add" {
        return n.doAdd(r)
    }

    return errors.New("Unknown command " + c)
}

func (*NetCmd) doInfo() error {
    fmt.Println("Net stats")
    return nil
}

func (*NetCmd) doAdd(args []string) error {
    if len(args) == 0 { return errors.New("Missing arguments for command 'net add'") }
    address := args[0]
    fmt.Printf("Connecting to %v \n", address)
    conn, err := net.Dial("tcp", address)
    if err == nil {
        go serv.Talk(conn)
    }
    return err
}
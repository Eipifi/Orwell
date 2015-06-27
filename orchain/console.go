package main
import (
    "bufio"
    "os"
    "fmt"
    "strings"
    "errors"
    "net"
)

type Command func([]string) error

var commands = map[string] Command {
    "connections": cmdConnections,
    "connect": cmdConnect,
}

func runConsole() {
    fmt.Println("Orchain v0.1")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf("> ")
        if ! scanner.Scan() { break }
        line := scanner.Text()
        if line == "x" { break }
        words := strings.Fields(line)
        if len(words) == 0 { continue }
        cmd, ok := commands[words[0]]
        if ok {
            if err := cmd(words[1:]); err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Printf("Unknown command %v\n", words[0])
        }
    }
}

// lists all connections
func cmdConnections([]string) error {
    fmt.Println("List all connections")
    return nil
}

// connects to a specified remote
func cmdConnect(args []string) error {
    if len(args) != 1 { return errors.New("Format: connect ADDRESS") }
    address := args[0]
    fmt.Printf("Connecting to address %v\n", address)
    c, err := net.Dial("tcp", address)
    if err != nil { return err }
    go Connect(c)
    return nil
}


/*

connect ADDRESS
    - connects to a specified remote

connections
    - lists all connections

buffer
    - lists all not yet confirmed transactions

mine ADDRESS
    - starts mining, with the specified address as destination

stats
    - number of blocks
    - head id
    - current difficulty
    - cumulative work done on the chain

balance ID
    - total sum of unspent bills
    - list of unspent bills for the key




*/
package main
import (
    "os"
    "fmt"
)

const usage = `Ortool is a handy tool for Orwell protocol.

Usage:

    ortool command [args]

The commands are:

    genkey      creates a new private key
    gencard     runs the interactive card generator
    read        reads the given PEM file and displays relevant info
    version     prints the ortool version
    config      prints the current configuration
    fetch       resolves the Orwell address and prints its value
    publish     publishes the card in the network
    chaininfo   displays current blockchain information

Use "ortool help [command]" for more information about a command.

`

func main() {
    // Skip the executable name
    args := os.Args[1:]
    // Command specified?
    if len(args) != 0 {
        command := commands[args[0]]
        if command == nil {
            fmt.Printf("Unknown command \"%s\".\n", args[0])
            os.Exit(2)
        }
        if err := command.Main(args[1:]); err != nil {
            if err == InvalidUsage {
                fmt.Println("Usage:", command.Usage())
                os.Exit(2)
            } else {
                fmt.Println("Error:", err)
                os.Exit(1)
            }
        }
    } else {
        fmt.Println(usage)
    }
}

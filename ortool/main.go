package main

import (
    "fmt"
    "os"
    "orwell/ortool/genkey"
    "orwell/ortool/gencard"
    "orwell/ortool/read"
    "orwell/ortool/fetch"
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
    chaininfo   displays current blockchain information

Use "ortool help [command]" for more information about a command.

`

func main() {

    if len(os.Args) == 1 {
        printUsage()
    }

    command := os.Args[1]
    args := os.Args[2:]

    if command == "help" {
        if len(args) == 0 {
            printUsage()
        } else if len(args) > 1 {
            fmt.Print("Too many arguments.\n\nUsage: ortool help <command>\n\n")
        } else {
            m, ok := modules[args[0]]
            if ok {
                fmt.Print(m.usage)
            } else {
                fmt.Print("No help entry for subcommand \"", args[0], "\".\n\n")
            }
        }
    } else {
        m, ok := modules[command]
        if ok {
            m.main(args)
        } else {
            fmt.Print("Unknown subcommand \"", command, "\"\nRun \"ortool help\" for usage.\n\n")
        }
    }
}

func printUsage() {
    fmt.Print(usage)
    os.Exit(0)
}

type module struct {
    usage string
    main func([]string)
}

var modules = map[string] module {
    "genkey": module{genkey.Usage, genkey.Main},
    "gencard": module{gencard.Usage, gencard.Main},
    "read": module{read.Usage, read.Main},
    "fetch": module{fetch.Usage, fetch.Main},
}
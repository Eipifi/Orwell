package main
import (
    "errors"
    "fmt"
)

type helpCommand struct {}

func (helpCommand) Usage() string {
    return "ortool help <command>"
}

func (helpCommand) Description() string {
    return ""
}

func (helpCommand) Main(args []string) error {
    if len(args) != 1 { return InvalidUsage }
    command := commands[args[0]]
    if command == nil {
        return errors.New("unknown command")
    } else {
        fmt.Println("Usage:", command.Usage(), "\n")
        fmt.Println(command.Description())
        return nil
    }
}

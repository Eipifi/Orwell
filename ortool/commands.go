package main
import "errors"

type Command interface {
    Usage() string
    Description() string
    Main([]string) error
}

var InvalidUsage = errors.New("Invalid usage")

var commands = map[string] Command {
    "help": &helpCommand{},
    "fetch": &fetchCommand{},
    "genkey": &genkeyCommand{},
    "gencard": &gencardCommand{},
    "publish": &publishCommand{},
    "read": &readCommand{},
}

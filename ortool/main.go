package main
import (
    "github.com/codegangsta/cli"
    "os"
    "orwell/ortool/cmd"
)

func main() {

    app := cli.NewApp()
    app.Name = "ortool"
    app.Usage = "a handy tool for Orwell protocol"
    app.Version = "0.0.1"
    app.EnableBashCompletion = true
    app.Author = "Michał Jabczyński"

    app.Commands = []cli.Command{
        {
            Name: "resolve",
            Aliases: []string{"r"},
            Usage: "resolve an Orwell address",
            Action: cmd.ResolveAction,
            Flags: []cli.Flag {

            },
        },
        {
            Name: "genkey",
            Aliases: []string{"g"},
            Usage: "generate a new key pair",
            Action: cmd.GenKeyAction,
            Flags: []cli.Flag {

            },
        },
        {
            Name: "gencard",
            Aliases: []string{"c"},
            Usage: "generate a new card",
            Action: cmd.GenCardAction,
            Flags: []cli.Flag {

            },
        },
    }

    app.Run(os.Args)
}

/*

- generate key
- generate card using the key

*/
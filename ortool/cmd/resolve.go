package cmd
import "github.com/codegangsta/cli"

func ResolveAction(c *cli.Context) {
    println("added task: ", c.Args().First())
}
package cmd
import (
    "bufio"
    "os"
    "fmt"
    "strings"
)

func Run(commands []Command) {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf("> ")
        if ! scanner.Scan() { break }
        line := scanner.Text()
        if line == "x" { break }
        words := strings.Fields(line)
        if len(words) == 0 { continue }

        for _, c := range commands {
            if c.Name() == words[0] {
                if err := c.Run(words[1:]); err != nil {
                    fmt.Println("Error: ", err)
                }
            }
        }
    }
}
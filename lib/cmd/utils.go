package cmd
import (
    "os"
    "bufio"
    "orwell/lib/foo"
    "fmt"
    "errors"
    "strings"
)

func PressEnterToContinue() {
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func ReadUint64(min, max uint64) (val uint64) {
    for {
        fmt.Printf("> ")
        _, err := fmt.Scanf("%d\n", &val)
        if err == nil {
            if val < min || val > max { err = errors.New("invalid number range") }
            if err == nil { break }
        }
        fmt.Print("Error: %v \n", err)
    }
    return
}

func ReadU256() (val foo.U256) {
    var err error
    for {
        hex_recipient := ReadString()
        val, err = foo.FromHex(strings.TrimSpace(hex_recipient))
        if err == nil { break }
        fmt.Printf("Error: %v \n", err)
    }
    return
}

func ReadString() string {
    for {
        fmt.Printf("> ")
        line, err := bufio.NewReader(os.Stdin).ReadString('\n')
        if err == nil {
            return trimSuffix(line, "\n")
        }
        fmt.Printf("Error: %v \n", err)
    }
}

func trimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}
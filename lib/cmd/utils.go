package cmd
import (
    "os"
    "bufio"
)

func PressEnterToContinue() {
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}
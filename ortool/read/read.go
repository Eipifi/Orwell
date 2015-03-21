package read
import (
    "fmt"
    "os"
    "io/ioutil"
    "encoding/pem"
    "orwell/orlib/sig"
)

const Usage = `usage: ortool read [path]

Reads the given file (card or key) and displays relevant info.

`
func Main(args []string) {

    var f *os.File
    var err error

    if len(args) > 1 {
        fail("Too may arguments.")
    }

    if len(args) == 1 {
        f, err = os.Open(args[0])
        if err != nil {
            fail("Invalid path.")
        }
    } else {
        f = os.Stdin
    }

    data, err := ioutil.ReadAll(f)

    if err != nil {
        fail("Failed to read input.")
    }

    b, _ := pem.Decode(data)

    if b == nil {
        fail("Failed to read PEM file.")
    }

    key, err := sig.ParsePrvKey(b.Bytes)

    if err == nil {
        fmt.Println("PRIVATE KEY")
        fmt.Printf("ID: %x\n", key.PublicPart().Id())
    }


}

func fail(msg string) {
    fmt.Println(msg)
    os.Exit(1)
}
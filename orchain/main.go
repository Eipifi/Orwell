package main
import (
    "orwell/lib/logging"
    "os"
    "fmt"
    "orwell/lib/blockstore"
    "strconv"
)

var Storage blockstore.BlockStorage

func main() {
    var port uint16 = 1984
    if len(os.Args) > 1 {
        v, err := strconv.Atoi(os.Args[1])
        if err == nil {
            port = uint16(v)
        } else {
            fmt.Println("Failed to parse port, using default")
        }
    }

    // Redirect logging to the file
    logFile, err := os.Create(fmt.Sprintf("/tmp/orchain.%v.log", port))
    if err != nil {
        fmt.Println(err)
        return
    }
    logging.DirectTo(logFile)

    // Open the storage
    db, err := blockstore.Open(fmt.Sprintf("/tmp/blockstore-%v", port))
    if err != nil {
        fmt.Println(err)
        return
    }
    Storage = blockstore.NewBlockStore(db)

    go runServer(port)
    runConsole()
}
package main
import (
    "orwell/lib/logging"
    "os"
    "fmt"
    "orwell/lib/blockstore"
    "strconv"
    "orwell/lib/miner"
    "orwell/lib/butils"
)

var Storage blockstore.BlockStorage
var MinerSup *miner.MiningSupervisor

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

    // Start the mining supervisor
    MinerSup = miner.NewSupervisor(butils.Uint256{}, Storage)

    go runServer(port)
    runConsole()
}
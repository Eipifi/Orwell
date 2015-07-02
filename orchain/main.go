package main
import (
    "orwell/lib/logging"
    "orwell/lib/blockstore"
    "orwell/lib/miner"
    "log"
)

var Config *ConfigManager
var Storage blockstore.BlockStorage
var MinerSup *miner.MiningSupervisor

func main() {
    if err := initialize(); err != nil {
        log.Panicln(err)
    }
    go runServer(Config.Port())
    runConsole()
}

func initialize() (err error) {
    Config, err = InitConfig()
    if err != nil { return }
    if err = logging.DirectToFile(Config.RelPath("orchain.log")); err != nil { return }
    db, err := blockstore.Open(Config.RelPath("db"))
    if err != nil { return }
    Storage = blockstore.NewBlockStore(db)
    MinerSup = miner.NewSupervisor(Config.MinerAddress(), Storage)
    return
}
package main
import (
    "orwell/lib/config"
    "orwell/lib/logging"
    "orwell/lib/cmd"
    "orwell/orchain/command"
    "orwell/orchain/serv"
    "orwell/lib/db"
)

func main() {
    // Initialize
    config.LoadDefault()
    logging.DirectToFile(config.Path("orchain.log"))

    // Load the block storage
    db.Initialize(config.Path("chain.bdb"))

    // Run server routines
    go serv.RunServer(config.GetInt("port"))

    // Run the console
    cmd.Run([]cmd.Command{
        &command.Stats{},
        &command.Miner{},
        &command.Net{},
        &command.Wallet{},
    })
}

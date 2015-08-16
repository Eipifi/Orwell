package main
import (
    "orwell/lib/config"
    "orwell/lib/logging"
    "orwell/orchain/serv"
    "orwell/lib/db"
    "orwell/lib/fcli"
    "orwell/orchain/cmd"
)

func main() {
    // Initialize
    config.LoadDefault()
    logging.DirectToFile(config.Path("orchain.log"))

    // Load the block storage
    db.Initialize(config.Path("chain.bdb"))

    // Run the managers
    serv.Bootstrap()

    // Run server routines
    go serv.RunServer(config.GetInt("port"))

    // Run the command-line finite state machine
    fsm := fcli.NewFSM("> ")

    fsm.On("main", "stats", cmd.StatsHandler)
    fsm.On("main", "s", cmd.StatsHandler)
    fsm.On("main", "balance $U256", cmd.BalanceHandler)
    fsm.On("main", "mine $U256", cmd.MinerHandler)
    fsm.On("main", "net", cmd.NetStatsHandler)
    fsm.On("main", "net add $str", cmd.NetAddHandler)
    fsm.On("main", "resolve $str", cmd.ResolveHandler)
    fsm.On("main", "wallet generate", cmd.WalletGenerateHandler)
    fsm.On("main", "wallet", cmd.WalletHandler)
    fsm.On("main", "send", cmd.SendHandler)
    fsm.On("main", "exit", fcli.ExitHandler)
    fsm.On("main", "x", fcli.ExitHandler)

    fsm.Run("main")
}
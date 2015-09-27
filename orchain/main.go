package main
import (
    "orwell/lib/config"
    "orwell/lib/logging"
    "orwell/orchain/serv"
    "orwell/lib/db"
    "orwell/lib/fcli"
    "orwell/orchain/cmd"
    "orwell/lib/utils"
    "orwell/lib/obp"
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
    go func() {
        utils.Ensure(obp.Serve(config.GetInt("port"), serv.Talk))
    }()

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
    fsm.On("main", "block $uint64", cmd.BlockByNumHandler)
    fsm.On("main", "pop", cmd.PopHandler)
    fsm.On("main", "exit", fcli.ExitHandler)
    fsm.On("main", "x", fcli.ExitHandler)

    fsm.Run("main")
}
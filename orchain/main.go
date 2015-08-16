package main
import (
    "orwell/lib/config"
    "orwell/lib/logging"
    "orwell/orchain/serv"
    "orwell/lib/db"
    "orwell/lib/fcli"
    "orwell/orchain/cmds"
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

    fsm.On("main", "stats", cmds.StatsHandler)
    fsm.On("main", "s", cmds.StatsHandler)
    fsm.On("main", "balance $U256", cmds.BalanceHandler)
    fsm.On("main", "mine $U256", cmds.MinerHandler)
    fsm.On("main", "net add $str", cmds.NetAddHandler)
    fsm.On("main", "resolve $str", cmds.ResolveHandler)
    fsm.On("main", "wallet generate", cmds.WalletGenerateHandler)
    fsm.On("main", "wallet", cmds.WalletHandler)
    fsm.On("main", "send", cmds.SendHandler)
    fsm.On("main", "exit", fcli.ExitHandler)
    fsm.On("main", "x", fcli.ExitHandler)

    fsm.Run("main")
}

/*
stats - general info about the blockchain, and pending transactions
exit - close the program
wallet - lists all owned private keys
wallet generate - creates a new private key
balance [ADDRESS] - display the available amount, unspent bills and pending transactions
mine [ADDRESS] - start mining coins <press any key to stop>

resolve [NAME] - resolve the domain name
pay - interactive prompt for making transactions

net - network info, connected peers
net add [IP] - attempt connecting with the
*/

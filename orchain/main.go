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

    // Run the managers
    serv.Bootstrap()

    // Run server routines
    go serv.RunServer(config.GetInt("port"))

    // Run the console
    cmd.Run([]cmd.Command{
        &command.StatsCmd{},
        &command.MinerCmd{},
        &command.NetCmd{},
        &command.WalletCmd{},
        &command.BalanceCmd{},
        &command.LogCmd{},
        &command.PayCmd{},
    })
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

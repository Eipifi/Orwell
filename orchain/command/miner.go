package command
import (
    "orwell/lib/foo"
    "errors"
    "orwell/orchain/miner"
    "fmt"
    "orwell/lib/cmd"
)

type MinerCmd struct {
    run bool
    miner *miner.SimpleMiner
}

func (c *MinerCmd) Name() string {
    return "mine"
}

func (c *MinerCmd) Run(args []string) error {

    if len(args) != 1 { return errors.New("Invalid usage") }

    wallet_id, err := foo.FromHex(args[0])
    if err != nil { return err }

    miner := miner.StartMiner(wallet_id)
    fmt.Println("Press ENTER/RETURN to continue")
    cmd.PressEnterToContinue()
    miner.Stop()
    fmt.Println("Stopped")

    return nil
}
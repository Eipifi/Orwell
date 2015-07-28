package command
import (
    "orwell/lib/foo"
    "errors"
    "orwell/orchain/miner"
)

type MinerCmd struct {
    run bool
    miner *miner.Miner
}

func (c *MinerCmd) Name() string {
    return "mine"
}

func (c *MinerCmd) Run(args []string) error {
    if len(args) != 1 { return errors.New("Invalid usage") }

    if c.miner == nil {
        c.miner = miner.NewMiner(foo.ZERO)
    }

    if args[0] == "start" {
        c.miner.Run(true)
    }

    if args[0] == "stop" {
        c.miner.Run(false)
    }

    return nil
}
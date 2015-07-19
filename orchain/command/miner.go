package command
import (
    "orwell/lib/foo"
    "errors"
    "orwell/orchain/miner"
)

type Miner struct {
    run bool
    miner *miner.Miner
}

func (m *Miner) Name() string {
    return "mine"
}

func (m *Miner) Run(args []string) error {
    if len(args) != 1 { return errors.New("Invalid usage") }

    if m.miner == nil {
        m.miner = miner.NewMiner(foo.ZERO)
    }

    if args[0] == "start" {
        m.miner.Run(true)
    }

    if args[0] == "stop" {
        m.miner.Run(false)
    }

    return nil
}
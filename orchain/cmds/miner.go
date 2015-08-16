package cmds
import (
    "orwell/lib/fcli"
    "orwell/lib/foo"
    "orwell/orchain/miner"
    "fmt"
)

func MinerHandler(id foo.U256) fcli.Result {

    fsm := fcli.NewFSM("> ")
    var m *miner.SimpleMiner

    fsm.On("start", "", func() fcli.Result {
        m = miner.StartMiner(id)
        fmt.Println("Press ENTER/RETURN to continue")
        return fcli.Next("end")
    })

    fsm.On("end", "$str", func(string) fcli.Result {
        m.Stop()
        fmt.Println("Stopped")
        return fcli.Exit(nil)
    })

    return fsm.Run("start")
}
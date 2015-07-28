package command
import (
    "fmt"
    "orwell/lib/db"
    "orwell/lib/protocol/orchain"
    "time"
)

type StatsCmd struct{}

func (*StatsCmd) Name() string {
    return "stats"
}

func (*StatsCmd) Run(args []string) error {

    db.Get().View(func(t *db.Tx){
        state := t.GetState()
        head := t.GetHeaderByID(state.Head)
        fmt.Printf("Last block: %v (#%v)\n", state.Head, state.Length)
        fmt.Printf("Total work: %v (%v)\n", state.Work, state.Work.Big())
        fmt.Printf("Difficulty: %v (%v)\n", head.Difficulty, head.Difficulty.Big())
        fmt.Printf("Timestamp : %v (%v)\n", head.Timestamp, time.Unix(int64(head.Timestamp), 0).Format(time.RFC1123))
        fmt.Printf("Reward    : %v \n", orchain.GetReward(state.Length))
    })

    return nil
}
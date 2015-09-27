package cmd
import (
    "orwell/lib/db"
    "fmt"
    "time"
    "orwell/lib/protocol/orchain"
    "orwell/lib/fcli"
)

func StatsHandler() fcli.Result {

    db.Get().View(func(t *db.Tx){
        state := t.GetState()
        head := t.GetHeaderByID(state.Head)
        fmt.Printf("Last block: %v (#%v)\n", state.Head, state.Length - 1)
        fmt.Printf("Total work: %v (%v)\n", state.Work, state.Work.Big())
        fmt.Printf("Difficulty: %v (%v)\n", head.Difficulty, head.Difficulty.Big())
        fmt.Printf("Timestamp : %v (%v)\n", head.Timestamp, time.Unix(int64(head.Timestamp), 0).Format(time.RFC1123))
        fmt.Printf("Reward    : %v \n", orchain.GetReward(state.Length))
    })

    return nil
}

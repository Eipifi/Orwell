package command
import (
    "fmt"
    "orwell/lib/db"
)

type Stats struct{}

func (*Stats) Name() string {
    return "stats"
}

func (*Stats) Run(args []string) error {

    if len(args) > 0 && args[0] == "pop" {
        db.Get().Pop()
    }

    state := db.Get().State()

    fmt.Printf("# of blocks: %v\n", state.Length)
    fmt.Printf("last block: %v\n", state.Head)
    fmt.Printf("Total work: %v\n", state.Work)
    fmt.Printf("difficulty: %v\n", db.Get().GetHeaderByID(state.Head).Difficulty)

    return nil
}
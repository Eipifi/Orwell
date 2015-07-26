package command
import (
    "fmt"
    "orwell/lib/db"
)

type Stats struct{}

func (*Stats) Name() string {
    return "s"
}

func (*Stats) Run(args []string) error {

    state := db.GetDB().State()

    fmt.Printf("# of blocks: %v\n", state.Length)
    fmt.Printf("last block: %v\n", state.Head)
    fmt.Printf("Total work: %v\n", state.Work)
    fmt.Printf("difficulty: %v\n", db.GetDB().Difficulty())

    return nil
}
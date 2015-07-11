package command
import (
    "fmt"
    "orwell/lib/blockstore"
)

type Stats struct{}

func (*Stats) Name() string {
    return "stats"
}

func (*Stats) Run([]string) error {

    s := blockstore.Get()

    fmt.Printf("# of blocks: %v\n", s.Length())
    fmt.Printf("last block: %v\n", s.Head())
    fmt.Printf("difficulty: %v\n", s.GetHeaderByID(s.Head()).Difficulty)
    return nil
}
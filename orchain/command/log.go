package command
import (
    "fmt"
    "orwell/lib/db"
)

type LogCmd struct{}

func (*LogCmd) Name() string {
    return "log"
}

func (*LogCmd) Run(args []string) error {

    db.Get().View(func(t *db.Tx){
        state := t.GetState()
        var i uint64
        for i = 1; i < 64; i += 1 {
            if i > state.Length { break }
            hd := t.GetHeaderByNum(state.Length - i)
            fmt.Printf("%v (#%v) \n", hd.ID(), state.Length - i)
        }
    })

    return nil
}
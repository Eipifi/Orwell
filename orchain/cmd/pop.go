package cmd
import (
    "orwell/lib/db"
    "orwell/lib/fcli"
)

func PopHandler() fcli.Result {

    db.Get().Update(func(t *db.Tx){
        t.PopBlock()
    })

    return nil
}
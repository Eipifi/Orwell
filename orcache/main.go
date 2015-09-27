package main
import (
    "orwell/orcache/serv"
    "orwell/lib/utils"
    "orwell/lib/obp"
)

func main() {
    utils.Ensure(obp.Serve(2000, serv.Talk))
}

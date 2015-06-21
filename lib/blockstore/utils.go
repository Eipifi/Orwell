package blockstore
import "errors"

func ensure(err error) {
    if err != nil {
        panic(err)
    }
}

func assert(condition bool) {
    if ! condition {
        panic(errors.New("Assertion failed"))
    }
}
package utils
import "errors"

func Ensure(err error) {
    if err != nil {
        panic(err)
    }
}

func Assert(condition bool) {
    if ! condition {
        panic(errors.New("Assertion failed"))
    }
}
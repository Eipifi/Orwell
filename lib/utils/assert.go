package utils
import (
    "errors"
    "bytes"
)

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

func Cat(slices ...[]byte) []byte {
    buf := bytes.Buffer{}
    for _, s := range slices {
        _, err := buf.Write(s)
        Ensure(err)
    }
    return buf.Bytes()
}
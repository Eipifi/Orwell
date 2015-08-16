package fcli
import (
    "orwell/lib/foo"
    "strconv"
)

type InputParser func(string) (interface{}, error)

//////////////////////////////////////////////////////////////////////////////

func StringParser(input string) (interface{}, error) {
    return input, nil
}

func Uint64Parser(input string) (interface{}, error) {
    return strconv.ParseUint(input, 10, 64)
}

func U256Parser(input string) (interface{}, error) {
    return foo.FromHex(input)
}
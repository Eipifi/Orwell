package fcli

type Result interface { }

type NextStateResult string
type ExitResult struct { err error }
type ErrorResult struct { err error }

func Next(state_name string) Result {
    return NextStateResult(state_name)
}

func Exit(err error) Result {
    return ExitResult{err}
}
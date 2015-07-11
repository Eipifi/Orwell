package cmd

type Command interface {
    Name() string
    Run([]string) error
}
package command

type Wallet struct{}

func (*Wallet) Name() string { return "wallet" }

func (*Wallet) Run(args []string) error {



    return nil
}
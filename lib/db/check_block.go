package db
import (
    "orwell/lib/protocol/orchain"
    "github.com/deckarep/golang-set"
    "errors"
    "orwell/lib/utils"
)

func CheckBlockMerkleRoot(t *Tx, b *orchain.Block) (err error) {
    return b.CheckMerkleRoot()
}

func CheckBlockHasAllDomains(t *Tx, b *orchain.Block) (err error) {
    domains := t.DomainsToRegister(b.Transactions)

    domains_required := mapset.NewSet()
    for _, d := range domains {
        utils.Assert(domains_required.Add(d))
    }
    domains_included := mapset.NewSet()
    for _, d := range b.Domains {
        utils.Assert(domains_included.Add(d))
    }

    if ! domains_required.Equal(domains_included) {
        return errors.New("The domain set of the block differs from the expected one.")
        // print which domains constitute the difference ?
    }
    return nil
}

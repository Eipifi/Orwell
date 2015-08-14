package db
import (
    "orwell/lib/protocol/orchain"
)

// General shape of all checks that need to be performed before a block is inserted.
// The function should return a clear and readable message of what went wrong, or nil on success.
// The function should check the value in context of the current database state.

// IMPORTANT: SCREW PERFORMANCE. IT'S NOT WORTH IT.
// This is an important part of code, and it must be as readable as conceivably possible.

type HeaderCheck func(*Tx, *State, *orchain.Header) error
type BlockCheck func(*Tx, *orchain.Block) error
type TxnCheck func(*Tx, *orchain.Transaction, bool) error
type TxnsCheck func(*Tx, []orchain.Transaction) error
type DomainCheck func(*Tx, *orchain.Domain) error
type DomainsCheck func(*Tx, []orchain.Domain) error

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Tx) ValidateNewBlock(block *orchain.Block) (err error) {
    s := t.GetState()
    for _, check := range header_checks {
        if err = check(t, s, &block.Header); err != nil { return }
    }
    for _, check := range block_checks {
        if err = check(t, block); err != nil { return }
    }
    for i, txn := range block.Transactions {
        for _, check := range txn_checks {
            if err = check(t, &txn, i == 0); err != nil { return }
        }
    }
    for _, check := range txns_checks {
        if err = check(t, block.Transactions); err != nil { return }
    }
    for _, domain := range block.Domains {
        for _, check := range domain_checks {
            if err = check(t, &domain); err != nil { return }
        }
    }
    for _, check := range domains_checks {
        if err = check(t, block.Domains); err != nil { return }
    }

    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Here are the lists of all checks performed before the block is saved in the database.

var header_checks []HeaderCheck = []HeaderCheck{
    CheckHeaderBasics,
    CheckHeaderDifficulty,
    CheckHeaderTimestamp,
}

var block_checks []BlockCheck = []BlockCheck{
    CheckBlockMerkleRoot,
}

var txn_checks []TxnCheck = []TxnCheck{
    CheckTxnIsNew,
    CheckTxnProof,
    CheckTxnInputsUnspent,
    CheckTxnNoDoubleSpend,
    CheckTxnBalance,
}

var txns_checks []TxnsCheck = []TxnsCheck{
    CheckTxnsNoDuplicateTxns,
    CheckTxnsNoDoubleSpend,
    CheckTxnsBalance,
}

var domain_checks []DomainCheck = []DomainCheck{}

var domains_checks []DomainsCheck = []DomainsCheck{}

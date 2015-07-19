package db
import (
    "orwell/lib/foo"
    "orwell/lib/protocol/orchain"
)

var GENESIS_DIFFICULTY foo.U256 = foo.ONE

func GenesisBlock() (b *orchain.Block) {
    b = &orchain.Block{}
    b.Header.Previous = foo.ZERO                    // 0
    b.Header.Timestamp = 1435519412                 // 2016/01/01 00:00:00 GMT // TODO: revert to proper date
    b.Header.Difficulty = GENESIS_DIFFICULTY        // Genesis block has 0 difficulty.
    b.Header.Nonce = 0                              // Any hash value meets the target(0), so nonce is set to 0.
    b.Transactions = []orchain.Transaction{
        orchain.Transaction{
            Outputs: []orchain.Bill{
                orchain.Bill{
                    Target: foo.ZERO,               // The first batch of coins is sent to nonexistent "0 wallet".
                    Value: orchain.GetReward(0),
                },
            },
        },
    }
    _ = b.ComputeMerkleRoot()

    // TODO: make a hash check

    return
}

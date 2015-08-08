package orchain
import (
    "orwell/lib/foo"
    "orwell/lib/utils"
)

var GENESIS_DIFFICULTY foo.U256 = foo.ZERO

func GenesisBlock() (b *Block) {
    b = &Block{}
    b.Header.Previous = foo.ZERO                    // 0
    b.Header.Timestamp = 1435519412                 // 2016/01/01 00:00:00 GMT // TODO: revert to proper date
    b.Header.Difficulty = GENESIS_DIFFICULTY        // Genesis block has 0 difficulty.
    b.Header.Nonce = 0                              // Any hash value meets the target(0), so nonce is set to 0.
    b.Transactions = []Transaction{
        Transaction{
            Outputs: []Bill{
                Bill{
                    Target: foo.ZERO,               // The first batch of coins is sent to nonexistent "0 wallet".
                    Value: GetReward(0),
                },
            },
            Label: "TODO: wymyślić fajny label",
        },
    }
    b.Domains = []Domain{}
    utils.Ensure(b.ComputeMerkleRoot())
    // TODO: make a hash check
    return
}
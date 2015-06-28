package blockstore
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/butils"
)

const GENESIS_DIFFICULTY uint16 = 0

func GenesisBlock() (b *orchain.Block) {
    b = &orchain.Block{}
    b.Header.Previous = butils.Uint256{}            // 0
    b.Header.Timestamp = 1435519410                 // 2016/01/01 00:00:00 GMT // TODO: revert to proper date
    b.Header.Difficulty = GENESIS_DIFFICULTY        // Genesis block has 0 difficulty.
    b.Header.Nonce = 0                              // Any hash value meets the target(0), so nonce is set to 0.
    b.Transactions = []orchain.Transaction{
        orchain.Transaction{
            Outputs: []orchain.Bill{
                orchain.Bill{
                    Target: butils.Uint256{},       // The first batch of coins is sent to nonexistent "0 wallet".
                    Value: 50,                      // TODO: update the generated value
                },
            },
        },
    }
    _ = b.ComputeMerkleRoot()

    // Sanity check - verify if the block hash is correct
    //bid := b.Header.ID().String()
    //if "c80edaaa0022d35b2f4f34a7197a2ea99991c68265ad4c1e7015e72d354f79d1" != bid {
    //    panic("unexpected genesis block hash [" + bid + "]")
    //}

    return
}

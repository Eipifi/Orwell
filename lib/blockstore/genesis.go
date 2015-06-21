package blockstore
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/butils"
)

func GenesisBlock() (b *orchain.Block) {
    b = &orchain.Block{}
    b.Header.Previous = butils.Uint256{}            // 0
    b.Header.Timestamp = 1451606400                 // 2016/01/01 00:00:00 GMT
    b.Header.Difficulty = 0                         // Genesis block has 0 difficulty.
    b.Header.Nonce = 0                              // Any hash value meets the target(0), so nonce is set to 0.
    b.Transactions = []orchain.Transaction{
        orchain.Transaction{
            TimeLock: 0,                            // This transaction must be spent
            Outputs: []orchain.Bill{
                orchain.Bill{
                    Target: butils.Uint256{},       // The first batch of coins is sent to nonexistent "0 wallet".
                    Value: 50,                      // TODO: update the generated value
                },
            },
        },
    }
    _ = b.ComputeMerkleRoot()
    bid := b.Header.ID().String()
    if "6f525d28bc52675c89fdc691f933ef17d26af8b34743d2abc2d3ab2271dd5e37" != bid {
        panic("unexpected genesis block hash [" + bid + "]")
    }
    return
}

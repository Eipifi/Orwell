package orchain
import "orwell/lib/crypto/sig"

type DomainData struct {
    DomainName string
    Key sig.PubKey
    LeaseTime uint64
    Signature sig.Signature
}

/*
    Transactions just as in BTC - inputs, outputs

    A transaction with NameHash is special - a miner can include only a specified number of those in a block.
    This way, only the highest-paid transactions are accepted. We call those *namehash transactions* (NHT).

    After an NHT is accepted into the blockchain, a countdown begins - during the next 144 blocks (~24h),
    the domain data struct must be placed in the blockchain.

        a) the hash must match the domain data
        b) the domain data must be correct

*/

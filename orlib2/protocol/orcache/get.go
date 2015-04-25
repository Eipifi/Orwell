package orcache
import (
    "orwell/orlib2/crypto/hash"
    "orwell/orlib2/protocol/common"
)

type Get struct {
    Token common.Token
    TTL common.TTL
    ID *hash.ID
    Version uint64
}
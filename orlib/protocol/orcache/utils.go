package orcache
import (
    "orwell/orlib/butils"
    "orwell/orlib/protocol/common"
)

type Tokener interface {
    GetToken() common.Token
}

type ChunkWithToken interface {
    butils.Chunk
    Tokener
}
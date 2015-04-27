package main
import (
    "orwell/orlib/butils"
    "orwell/orlib/protocol/orcache"
)

func (p *Peer) handleMessage(msg butils.Chunk) {
    switch msg := msg.(type) {
        case *orcache.Get: go p.handleGet(msg)
        default: panic("Unrecognized message type")
    }
}

func (p *Peer) handleGet(msg *orcache.Get) {
    res := Find(msg) // recover from panic if channel is closed
    if res.Card != nil {
        p.Out <- &orcache.CardFound{msg.Token, res.Card}
    } else {
        p.Out <- &orcache.CardNotFound{msg.Token, res.TTL}
    }
}
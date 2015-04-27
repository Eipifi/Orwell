package main
import (
    "orwell/orlib/butils"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/crypto/hash"
)

func (p *Peer) handleMessage(msg butils.Chunk) {
    switch msg := msg.(type) {
        case *orcache.Get:              go p.handleGet(msg)
        case *orcache.CardFound:        p.handleCardFound(msg)
        case *orcache.CardNotFound:     p.handleCardNotFound(msg)
        default: panic("Unrecognized message type")
    }
}

func (p *Peer) handleGet(msg *orcache.Get) {
    p.Log.Println("Received Get for id:", msg.ID, "version:", msg.Version)
    res := Find(msg) // recover from panic if channel is closed
    if res.Card != nil {
        p.Log.Println("Sending CardFound for id:", msg.ID, "version:", msg.Version)
        p.Out <- &orcache.CardFound{msg.Token, res.Card}
    } else {
        p.Log.Println("Sending CardNotFound for id:", msg.ID, "version:", msg.Version)
        p.Out <- &orcache.CardNotFound{msg.Token, res.TTL}
    }
}

func (p *Peer) handleCardFound(msg *orcache.CardFound) {
    order, ok := p.Orders[msg.Token]
    if ok {
        // FIXME: Do something about the int64/uint64 incompatibility issue
        if hash.Equal(msg.Card.Key.Id(), order.Req.ID) && uint64(msg.Card.Payload.Version) >= order.Req.Version {
            go func(){ order.Sink <- &GetResult{msg.Card, 0} }()
            delete(p.Orders, msg.Token)
            p.Log.Println("Received CardFound for id:", order.Req.ID, "version:", order.Req.Version)
            return
        } else {
            p.Log.Println("Received Cardfound card does not match the requirements")
        }
    } else {
        p.Log.Println("Received CardFound has an invalid token")
    }
    p.close()
}

func (p *Peer) handleCardNotFound(msg *orcache.CardNotFound) {
    order, ok := p.Orders[msg.Token]
    if ok {
        go func(){ order.Sink <- &GetResult{nil, msg.TTL} }()
        delete(p.Orders, msg.Token)
        p.Log.Println("Received CardNotFound for id:", order.Req.ID, "version:", order.Req.Version)
        return
    }
    p.Log.Println("Received CardNotFound has an invalid token")
    p.close()
}
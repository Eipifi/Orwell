package main
import (
    "orwell/orlib/butils"
    "orwell/orlib/protocol/orcache"
)

func (p *Peer) handleMessage(msg butils.Chunk) {
    switch msg := msg.(type) {
        case *orcache.GetReq:       p.handleGetReq(msg)
        case *orcache.GetRsp:       p.handleGetRsp(msg)
        case *orcache.PublishReq:   p.handlePublishReq(msg)
        case *orcache.PublishRsp:   p.handlePublishRsp(msg)
        default: panic("Unrecognized message type")
    }
}

func (p *Peer) handleGetReq(msg *orcache.GetReq) {
    go p.AsyncSend(Find(msg))
}

func (p *Peer) handlePublishReq(msg *orcache.PublishReq) {
    go p.AsyncSend(Publish(msg))
}

func (p *Peer) handleGetRsp(msg *orcache.GetRsp) {
    if !p.GetOrders.Respond(msg) { p.close() }
}

func (p *Peer) handlePublishRsp(msg *orcache.PublishRsp) {
    if !p.PutOrders.Respond(msg) { p.close() }
}
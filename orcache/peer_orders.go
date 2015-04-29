package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/protocol/common"
)

type GetOrder struct {
    Req *orcache.Get
    Sink chan *GetResult
}

type PutOrder struct {
    Req *orcache.Publish
    Sink chan *common.TTL
}

func (p *Peer) AskGet(req *orcache.Get) (result *GetResult) { // Called outside of peer.Lifetime()
    defer recover()
    sink := make(chan *GetResult)
    p.ToGet <- GetOrder{req, sink}
    result = <- sink
    return
}

func (p *Peer) handleGetOrder(order GetOrder) {
    if _, ok := p.GetOrders[order.Req.Token]; ok {
        order.Sink <- nil // order token already in use
        return
    }
    p.GetOrders[order.Req.Token] = order
    p.AsyncSend(order.Req)
}

func (p *Peer) AskPut(req *orcache.Publish) (result *common.TTL) {
    defer recover()
    sink := make(chan *common.TTL)
    p.ToPut <- PutOrder{req, sink}
    result = <- sink
    return
}

func (p *Peer) handlePutOrder(order PutOrder) {
    if _, ok := p.PutOrders[order.Req.Token]; ok {
        order.Sink <- nil // order token already in use
        return
    }
    p.PutOrders[order.Req.Token] = order
    p.AsyncSend(order.Req)
}

package main
import "orwell/orlib/protocol/orcache"

type GetOrder struct {
    Req *orcache.Get
    Sink chan *GetResult
}

func (p *Peer) AskGet(req *orcache.Get) (result *GetResult) { // Called outside of peer.Lifetime()
    defer recover()
    sink := make(chan *GetResult)
    p.ToDo <- GetOrder{req, sink}
    sink <- result
    return
}

func (p *Peer) handleOrder(order GetOrder) {
    if _, ok := p.Orders[order.Req.Token]; ok {
        order.Sink <- nil // order token already in use
        return
    }
    p.Orders[order.Req.Token] = order
    p.Out <- order.Req
}
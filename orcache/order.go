package main
import (
    "orwell/orlib/protocol/common"
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/butils"
)

type request struct {
    msg orcache.ChunkWithToken
    validator func(orcache.ChunkWithToken) bool
    sink chan orcache.ChunkWithToken
}

type response struct {
    msg orcache.ChunkWithToken
    sink chan bool
}

func (o *request) cancel() {
    o.complete(nil)
}

func (o *request) complete(rsp orcache.ChunkWithToken) {
    o.sink <- rsp
}

type OrderManager struct {
    sink chan<- butils.Chunk
    orders map[common.Token] *request
    ords chan *request
    rsps chan *response
}

func NewOrderManager(sink chan<- butils.Chunk) *OrderManager {
    m := &OrderManager{}
    m.sink = sink
    m.orders = make(map[common.Token] *request)
    m.ords = make(chan *request)
    m.rsps = make(chan *response)
    go m.lifecycle()
    return m
}

func (m *OrderManager) Ask(msg orcache.ChunkWithToken, validator func(orcache.ChunkWithToken) bool) orcache.ChunkWithToken {
    defer recover()
    ord := &request{msg, validator, make(chan orcache.ChunkWithToken)}
    m.ords <- ord
    result := <- ord.sink
    return result
}

func (m *OrderManager) Respond(msg orcache.ChunkWithToken) bool {
    defer recover()
    rsp := &response{msg, make(chan bool)}
    m.rsps <- rsp
    result := <- rsp.sink
    return result
}

func (m *OrderManager) Close() {
    // TODO: kill the lifecycle and close channels
}

func (m *OrderManager) handleOrder(req *request) {
    if _, ok := m.orders[req.msg.GetToken()]; ok {
        req.cancel()
    } else {
        m.orders[req.msg.GetToken()] = req
        m.sink <- req.msg
    }
}

func (m *OrderManager) handleResponse(rsp *response) {
    if req, ok := m.orders[rsp.msg.GetToken()]; ok {
        delete(m.orders, rsp.msg.GetToken())
        if req.validator == nil || req.validator(rsp.msg) {
            rsp.sink <- true
            req.complete(rsp.msg)
        } else {
            rsp.sink <- false
            req.cancel()
        }
    } else {
        rsp.sink <- false
    }
}

func (m *OrderManager) lifecycle() {
    for {
        select {
            case req := <- m.ords: m.handleOrder(req)
            case rsp := <- m.rsps: m.handleResponse(rsp)
        }
    }
}
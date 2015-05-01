package main
import (
    "orwell/orlib/protocol/orcache"
    "orwell/orlib/butils"
    "orwell/orlib/protocol/common"
    "sync"
)

type request struct {
    msg orcache.ChunkWithToken
    validator func(orcache.ChunkWithToken) bool
    sink chan orcache.ChunkWithToken
}

func (r *request) complete(msg orcache.ChunkWithToken) {
    r.sink <- msg
}

func (r *request) cancel() {
    r.complete(nil)
}

type response struct {
    msg orcache.ChunkWithToken
    sink chan bool
}

type RequestRouter struct {
    sink chan<- butils.Chunk
    orders map[common.Token] *request
    mtx *sync.Mutex
    closed bool
}

func NewRouter(sink chan<- butils.Chunk) *RequestRouter {
    m := &RequestRouter{}
    m.sink = sink
    m.orders = make(map[common.Token] *request)
    m.mtx = &sync.Mutex{}
    m.closed = false
    return m
}

func (m *RequestRouter) Ask(msg orcache.ChunkWithToken, validator func(orcache.ChunkWithToken) bool) orcache.ChunkWithToken {
    ord := &request{msg, validator, make(chan orcache.ChunkWithToken)}
    go m.handleOrder(ord)
    return <- ord.sink
}

func (m *RequestRouter) Respond(msg orcache.ChunkWithToken) bool {
    rsp := &response{msg, make(chan bool)}
    go m.handleResponse(rsp)
    return <- rsp.sink
}

func (m *RequestRouter) Close() {
    m.mtx.Lock()
    if m.closed {
        m.mtx.Unlock()
    } else {
        m.closed = true
        m.mtx.Unlock()
        for _, v := range m.orders {
            go v.cancel()
        }
    }
}

func (m *RequestRouter) handleOrder(req *request) {
    m.mtx.Lock()
    if m.closed {
        m.mtx.Unlock()
        req.cancel()
    } else {
        if _, ok := m.orders[req.msg.GetToken()]; ok {
            m.mtx.Unlock()
            req.cancel()
        } else {
            m.orders[req.msg.GetToken()] = req
            m.mtx.Unlock()
            m.sink <- req.msg
        }
    }
}

func (m *RequestRouter) handleResponse(rsp *response) {
    m.mtx.Lock()
    if m.closed {
        m.mtx.Unlock()
        rsp.sink <- false
    } else {
        if req, ok := m.orders[rsp.msg.GetToken()]; ok {
            delete(m.orders, rsp.msg.GetToken())
            m.mtx.Unlock()
            if req.validator == nil || req.validator(rsp.msg) {
                rsp.sink <- true
                req.complete(rsp.msg)
            } else {
                rsp.sink <- false
                req.cancel()
            }
        } else {
            m.mtx.Unlock()
            rsp.sink <- false
        }
    }
}
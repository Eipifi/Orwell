package serv
import (
    "log"
    "orwell/lib/obp"
    "net"
    "orwell/lib/logging"
    "orwell/lib/protocol/orchain"
    "errors"
    "orwell/lib/foo"
    "orwell/lib/utils"
    "orwell/lib/db"
)

type Peer struct {
    conn *obp.MsgConn
    log *log.Logger
}

func Talk(socket net.Conn) {
    if err := TalkTo(socket); err != nil {
        log.Println(err)
    }
}

func TalkTo(socket net.Conn) (err error) {
    p := &Peer{}
    p.log = logging.GetLogger(socket.RemoteAddr().String())
    p.conn = orchain.Connection(socket)

    defer p.conn.Close()
    ConnMgr().Join(p)
    for {
        err = p.conn.Handle(p.messageHandler)
        if err != nil {
            break
        }
    }
    ConnMgr().Leave(p)
    return nil
}

func (p *Peer) Close() {
    p.conn.Close()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *Peer) messageHandler(msg obp.Msg) (rsp obp.Msg, err error) {

    switch req := msg.(type) {
        case *orchain.MsgHead:      return p.handleMsgHead(req)
        case *orchain.MsgGetBlock:  return p.handleMsgGetBlock(req)
        case *orchain.MsgGetTxns:   return p.handleGetTxns(req)
    }

    return nil, errors.New("Unknown message type")
}

func (p *Peer) handleMsgHead(req *orchain.MsgHead) (rsp *orchain.MsgTail, err error) {
    //    Possible responses:
    //        - Same state, nothing to do                             // eq work, no headers
    //        - I have more work, these are the blocks you're missing // hi work, headers
    //        - I have more work, but I do not know your head block   // hi work, no headers
    //        - I have less, can't help ya                            // lo work, no headers

    state := db.Get().State()

    rsp = &orchain.MsgTail{}
    rsp.Work = state.Work

    cmp := foo.Compare(req.Work, rsp.Work)

    if cmp >= 0 {
        // remote side has more work, we can't help
        return
    } else {
        // we have more - wonder if we can help here
        num_ptr := db.Get().GetNumByID(req.Id)

        // If the header is not known, we can't send any subsequent headers
        if num_ptr == nil { return }

        // But if we already know this header, we can help - let's send the rest
        for num := 1 + *num_ptr; num < state.Length; num += 1 {
            header := db.Get().GetHeaderByNum(num)
            if header == nil { break }
            rsp.Headers = append(rsp.Headers, *header)
        }
    }
    return
}

func (p *Peer) handleMsgGetBlock(req *orchain.MsgGetBlock) (rsp *orchain.MsgBlock, err error) {
    return &orchain.MsgBlock{db.Get().GetBlockByID(req.ID)}, nil
}

func (p *Peer) handleGetTxns(req *orchain.MsgGetTxns) (rsp *orchain.MsgTxns, err error) {
    rsp = &orchain.MsgTxns{}
    db.Get().View(func(t *db.Tx) {
        rsp.Transactions = t.UnconfirmedTransactions()
    })
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var ErrInvalidResponse = errors.New("Invalid response type")

func (p *Peer) AskHead(revert uint64) (*orchain.MsgTail, error) {
    utils.Assert(revert >= 1)

    state := db.Get().State()
    req := &orchain.MsgHead{}
    req.Work = state.Work

    num := state.Length
    if revert > num {
        num = 0
    } else {
        num -= revert
    }
    id := db.Get().GetIDByNum(num)
    utils.Assert(id != nil)
    req.Id = *id

    rsp, err := p.conn.Query(req)
    if err != nil { return nil, err }
    rsp_cast, ok := rsp.(*orchain.MsgTail)
    if ok { return rsp_cast, nil }
    return nil, ErrInvalidResponse
}

func (p *Peer) AskBlock(id foo.U256) (*orchain.MsgBlock, error) {
    rsp, err := p.conn.Query(&orchain.MsgGetBlock{id})
    if err != nil { return nil, err }
    rsp_cast, ok := rsp.(*orchain.MsgBlock)
    if ok { return rsp_cast, nil }
    return nil, ErrInvalidResponse
}

func (p *Peer) AskTxns() (*orchain.MsgTxns, error) {
    rsp, err := p.conn.Query(&orchain.MsgGetTxns{})
    if err != nil { return nil, err }
    rsp_cast,  ok := rsp.(*orchain.MsgTxns)
    if ok { return rsp_cast, nil }
    return nil, ErrInvalidResponse
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

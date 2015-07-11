package serv
import (
    "log"
    "orwell/lib/obp"
    "net"
    "orwell/lib/logging"
    "orwell/lib/protocol/orchain"
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

    for {
        err = p.conn.Handle(p.messageHandler)
        if err != nil {
            return
        }
    }
}

func (p *Peer) messageHandler(msg obp.Msg) (rsp obp.Msg, err error) {

    /*
    switch req := msg.(type) {
        case *orchain.MsgHead:

            //    Possible responses:
            //        - Same state, nothing to do                             // eq work, no headers
            //        - I have more work, these are the block you're missing  // hi work, headers
            //        - I have more work, but I do not know your head block   // hi work, no headers
            //        - I have less, can't help ya                            // lower work, no headers

            response := &orchain.MsgTail{}
            response.Work = blockstore.Get().Work()
    }
    */
    return nil, nil
}
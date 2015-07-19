package serv
import (
    "log"
    "orwell/lib/logging"
    "math/rand"
    "time"
    "orwell/lib/protocol/orchain"
    "orwell/lib/db"
    "orwell/lib/foo"
    "errors"
    "github.com/deckarep/golang-set"
)

type Manager struct {
    log *log.Logger
    peers mapset.Set
}

func (m *Manager) Join(peer *Peer) {
    m.log.Printf("Peer joined: %v", peer)
    m.peers.Add(peer)
}

func (m *Manager) Leave(peer *Peer) {
    m.log.Printf("Peer left: %v", peer)
    m.peers.Remove(peer)
}

func (m *Manager) GetRandomPeers(num int) []*Peer {
    peers := m.peers.ToSlice()
    if num > len(peers) {
        num = len(peers)
    }
    result := make([]*Peer, num)
    // TODO randomize order
    for i := 0; i < num; i += 1 {
        result[i], _ = peers[i].(*Peer)
    }
    return result
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var instance *Manager

func makeInstance() *Manager {
    mgr := &Manager{}
    mgr.log = logging.GetLogger("")
    mgr.peers = mapset.NewSet()
    go mgr.syncLoop()
    return mgr
}

func GetManager() *Manager {
    if instance == nil { instance = makeInstance() }
    return instance
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *Manager) syncLoop() {
    for {
        peers := m.GetRandomPeers(1)
        if len(peers) > 0 {
            if err := m.sync(peers[0]); err != nil {
                peers[0].Close()
            }
        }
        interval := 3 + rand.Intn(3)
        time.Sleep(time.Duration(interval) * time.Second)
    }
}

func (m *Manager) sync(peer *Peer) (err error) {
    rsp := &orchain.MsgTail{}
    state := db.Get().State()
    var revert uint64 = 1
    for len(rsp.Headers) == 0 {
        if rsp, err = peer.AskHead(revert); err != nil { return }
        if foo.Compare(rsp.Work, state.Work) != 1 { return nil }
        if revert > state.Length && len(rsp.Headers) == 0 { return errors.New("Node advertises more work, yet does not send any headers after genesis block") }
        revert *= 2
    }
    // TODO verify if headers are properly signed, sum up to the declared work, and make overall sense

    // Iterate to first unknown header
    headers := rsp.Headers
    for {
        if len(headers) == 0 { return }
        if db.Get().GetNumByID(headers[0].ID()) == nil { break }
        headers = headers[1:]
    }

    // Drop the obsolete blocks
    for (headers[0].Previous) != (db.Get().State().Head) {
        db.Get().Pop()
    }

    // Download blocks and apply in order
    for _, h := range headers {
        var block_rsp *orchain.MsgBlock
        if block_rsp, err = peer.AskBlock(h.ID()); err != nil { return }
        if block_rsp.Block == nil { return } // The peer promised to deliver the block, and failed - what to do?
        if err = db.Get().Push(block_rsp.Block); err != nil { return }
    }

    m.log.Printf("Sync successful, downloaded %v blocks", len(headers))
    return
}
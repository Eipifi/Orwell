package serv
import (
    "log"
    "github.com/deckarep/golang-set"
    "orwell/lib/logging"
)

type ConnManager struct {
    log *log.Logger
    peers mapset.Set // this set is thread-safe!
}

func (m *ConnManager) Join(peer *Peer) {
    m.log.Printf("Peer joined: %v", peer)
    m.peers.Add(peer)
}

func (m *ConnManager) Leave(peer *Peer) {
    m.log.Printf("Peer left: %v", peer)
    m.peers.Remove(peer)
}

func (m *ConnManager) GetAllPeers() (peers []*Peer) {
    for _, p := range m.peers.ToSlice() {
        peers = append(peers, p.(*Peer))
    }
    return
}

func (m *ConnManager) GetRandomPeers(num int) []*Peer {
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

var connInstance *ConnManager

func ConnMgr() *ConnManager { // TODO: synchronize
    if connInstance == nil {
        connInstance = &ConnManager{}
        connInstance.log = logging.GetLogger("")
        connInstance.peers = mapset.NewSet()
    }
    return connInstance
}
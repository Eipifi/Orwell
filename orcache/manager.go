package main
import (
    "log"
    "os"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
    "stathat.com/c/jconfig"
    "fmt"
    "orwell/orlib/protocol/orcache"
)

type PeerFinder interface {
    FindPeer(hash.ID) *Peer
}

type Manager struct {
    log *log.Logger
    AdvertisedPort common.Port
    ActualPort common.Port
    AdvertisedID *hash.ID
    cfg *jconfig.Config
}

func NewManager() *Manager {
    m := &Manager{}
    m.log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    m.cfg = jconfig.LoadConfig("/Users/eipifi/Go/src/orwell/config/default.orcache.json") // DEVELOPMENT
    return m
}

func (m *Manager) Join(peer *Peer) {
    m.log.Println("Joined", peer.Hs)
}

func (m *Manager) Leave(peer *Peer) {
    m.log.Println("Left", peer.Hs)
}

func (m *Manager) FindPeer(id hash.ID) *Peer {
    return nil
}

func (m *Manager) CheckHandshake(hs *orcache.Handshake) bool {
    return true
}

func (m *Manager) FindAddresses(*hash.ID) []common.Address {
    return nil
}

func main() {
    m := NewManager()
    fmt.Println(m.Run())
}
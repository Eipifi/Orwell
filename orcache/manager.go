package main
import (
    "orwell/orlib/crypto/hash"
    "log"
    "os"
)

var Manager *PeerManager = NewManager()

func NewManager() *PeerManager {
    return &PeerManager{Log: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)}
}

type PeerManager struct {
    Log *log.Logger
}

func (m *PeerManager) Join(p *Peer) {
    m.Log.Println("Peer joined:", p.Hs)
}

func (m *PeerManager) Leave(p *Peer) {
    m.Log.Println("Peer left:", p.Hs)
}

func (m *PeerManager) FindPeer(id *hash.ID) *Peer {
    return nil
}

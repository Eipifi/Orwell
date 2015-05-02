package main
import (
    "log"
    "os"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
)

type PeerFinder interface {
    Find(*hash.ID) *Peer
}

type Manager interface {
    PeerFinder
    Join(*Peer)
    Leave(*Peer)
    LocalAddress() *common.Address
}

type ManagerImpl struct {
    log *log.Logger
    address *common.Address
}

func NewManagerImpl() Manager {
    m := &ManagerImpl{}
    m.log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    return m
}

func (m *ManagerImpl) Join(peer *Peer) {
    m.log.Println("Joined:", peer.Hs)
}

func (m *ManagerImpl) Leave(peer *Peer) {
    m.log.Println("Left:", peer.Hs)
}

func (m *ManagerImpl) Find(*hash.ID) *Peer {
    return nil
}

func (m *ManagerImpl) LocalAddress() *common.Address {
    return m.address
}
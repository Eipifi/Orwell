package main
import (
    "orwell/orlib/crypto/hash"
    "sort"
)

type RingMap struct {
    peers map[hash.ID] *Peer
    ring hash.IdList
}

func NewRingMap() *RingMap {
    m := &RingMap{}
    m.peers = make(map[hash.ID] *Peer)
    return m
}

func (m *RingMap) Put(peer *Peer) {
    id := peer.Hs.Address.Id()
    if m.ring.Contains(id) {

    }
    m.ring = append(m.ring, id)
    m.peers[id] = peer
    sort.Sort(m.ring)
}

func (m *RingMap) Contains(id hash.ID) bool {
    return m.ring.Contains(id)
}

func (m *RingMap) Del(peer *Peer) {
    id := peer.Hs.Address.Id()
    if m.Contains(id) {
        n := m.ring.ClosestNotSmaller(id)
        m.ring = append(m.ring[:n], m.ring[n+1:]...)
        delete(m.peers, id)
    }
}

func (m *RingMap) Nearest(id hash.ID, n int) []*Peer {
    ids := m.ring.NClosest(id, n)
    res := make([]*Peer, len(ids))
    for i := 0; i < len(ids); i++ {
        res[i] = m.peers[ids[i]]
    }
    return res
}
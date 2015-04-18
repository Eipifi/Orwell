package main
import (
    "orwell/orlib/sig"
    "net"
)

type Manager struct {
    socket net.Listener
}

func NewManager(connectionString string) (m *Manager, err error) {
    m = &Manager{}
    m.socket, err = net.Listen("tcp", connectionString)
    return
}

func (m *Manager) Lifecycle() {
    env := &Env{m, &EmptyCache{}, &EmptyTokenLocker{}}
    defer m.socket.Close()
    for {
        conn, err := m.socket.Accept()
        if err != nil { break }
        peer := NewPeer(conn, env)
        go peer.Lifecycle()
    }
}

func (m *Manager) PickPeer(id sig.ID) *Peer {
    return nil
}
package main
import (
    "net"
    "orwell/orlib/protocol/types"
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
        peer := NewPeer(env)
        go peer.Lifecycle(conn)
    }
}

func (m *Manager) PickPeer(id *types.ID) *Peer {
    return nil
}
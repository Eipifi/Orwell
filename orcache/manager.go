package main
import (
    "log"
    "os"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
    "net"
    "stathat.com/c/jconfig"
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
    cfg *jconfig.Config
}

func NewManagerImpl() *ManagerImpl {
    m := &ManagerImpl{}
    m.log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    m.cfg = jconfig.LoadConfig("/Users/eipifi/Go/src/orwell/config/default.orcache.json") // DEVELOPMENT
    return m
}

func (m *ManagerImpl) parseAddress() {

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

func (m *ManagerImpl) Run() error {
    socket, err := net.Listen("tcp", ":" + m.cfg.GetString("port"))
    if err != nil { return err }
    for {
        conn, err := socket.Accept()
        if err != nil { return err }
        NewPeer(conn, m)
    }
}

func main() {
    m := NewManagerImpl()
    m.Run()
}
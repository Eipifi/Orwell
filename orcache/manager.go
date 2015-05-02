package main
import (
    "log"
    "os"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/protocol/common"
    "net"
    "stathat.com/c/jconfig"
    "fmt"
    "strconv"
    "orwell/orlib/netutils"
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
    m.setAddress()
    m.log.Println("Started orcache with address", m.address)
    return m
}

func (m *ManagerImpl) setAddress() {
    m.address = &common.Address{}
    m.address.Nonce = uint64(m.cfg.GetInt("nonce"))
    m.address.Port = uint16(m.cfg.GetInt("port"))
    if m.cfg.GetString("ip") == "" {
        if m.address.IP = netutils.FindExternalIp(); m.address.IP == nil {
            panic("External IP not found")
        }
    } else {
        m.address.IP = net.ParseIP(m.cfg.GetString("ip"))
    }
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
    socket, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(m.address.Port), 10))
    if err != nil { return err }
    for {
        conn, err := socket.Accept()
        if err != nil { return err }
        NewPeer(conn, m)
    }
}

func main() {
    m := NewManagerImpl()
    fmt.Println(m.Run())
}
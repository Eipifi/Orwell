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

func (m *ManagerImpl) Run() error {
    // Setup the address
    m.address = &common.Address{}
    m.address.Nonce = uint64(m.cfg.GetInt("nonce"))
    if m.cfg.GetString("ip") == "" {
        if m.address.IP = netutils.FindExternalIp(); m.address.IP == nil {
            panic("External IP not found")
        }
    } else {
        m.address.IP = net.ParseIP(m.cfg.GetString("ip"))
    }
    if m.cfg.GetBool("port_is_external") {
        m.address.Port = uint16(m.cfg.GetInt("port"))
    } else {
        m.address.Port = 0
    }
    // Setup the port
    var actual_port uint16 = uint16(m.cfg.GetInt("port"))
    if m.cfg.GetBool("try_upnp") {
        if upnp_port := netutils.FindExternalUpnpPort(); upnp_port != 0 {
            actual_port = upnp_port
        } else {
            m.log.Println("Failed to negotiate an external UPNP port, using default")
        }
    }
    var ext_msg string = ""
    if m.address.IsInternal() {
        ext_msg = "(port 0 means an externally unreachable NAT address)"
    }
    m.log.Println("Advertised address", m.address, ext_msg)
    socket, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(actual_port), 10))
    m.log.Println("Listening on TCP port", actual_port)
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
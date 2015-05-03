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
    FindPeer(*hash.ID) *Peer
}

type Manager interface {
    PeerFinder
    Join(*Peer)
    Leave(*Peer)
    LocalAddress() *common.Address
    FindAddresses(*hash.ID) []common.Address
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

func (m *ManagerImpl) FindPeer(*hash.ID) *Peer {
    return nil
}

func (m *ManagerImpl) FindAddresses(*hash.ID) []common.Address {
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
        // If ip not specified in the configuration
        if m.address.IP = netutils.FindExternalIp(); m.address.IP == nil {
            // Maybe use local ip upon first established connection?
            m.log.Fatalln("Failed to fetch external IP address")
        }
    } else {
        m.address.IP = net.ParseIP(m.cfg.GetString("ip"))
    }

    // Setup the port
    m.address.Port = uint16(m.cfg.GetInt("port"))
    actual_port := m.address.Port
    if ! m.cfg.GetBool("port_is_external") {
        m.address.Port = 0
        if m.cfg.GetBool("try_upnp") {
            if upnp_port := netutils.FindExternalUpnpPort(); upnp_port != 0 {
                actual_port = upnp_port
                m.address.Port = upnp_port
            } else {
                m.log.Println("Failed to negotiate an external UPNP port, using default")
            }
        }
    }

    // Connect
    m.log.Println("Advertised address:", m.address)
    m.log.Println("Advertised id:", m.address.Id())
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
package main
import (
    "orwell/orlib/protocol/common"
    "orwell/orlib/netutils"
    "net"
    "strconv"
    "orwell/orlib/crypto/hash"
)

func (m *Manager) Run() error {
    // Setup the ID
    id_str := m.cfg.GetString("id")
    var id_val hash.ID
    if id_str == "" {
        id_val = hash.NewRandomId()
    } else {
        var err error
        id_val, err = hash.HexToID(id_str)
        if err != nil {
            panic(err)
        }
    }
    m.AdvertisedID = &id_val

    // Setup the port
    m.AdvertisedPort = common.Port(m.cfg.GetInt("port"))
    m.ActualPort = m.AdvertisedPort
    if ! m.cfg.GetBool("port_is_external") {
        if m.cfg.GetBool("try_upnp") {
            if upnp_port := netutils.FindExternalUpnpPort(m.ActualPort); upnp_port != 0 {
                m.AdvertisedPort = upnp_port
            } else {
                m.log.Println("Failed to negotiate an external UPNP port, using default")
                m.AdvertisedPort = common.NoPort
            }
        }
    }

    // Connect
    m.log.Println("Advertised port:", m.AdvertisedPort)
    m.log.Println("Advertised id:", m.AdvertisedID)
    socket, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(m.ActualPort), 10))
    m.log.Println("Listening on TCP port", m.ActualPort)
    if err != nil { return err }
    for {
        conn, err := socket.Accept()
        if err != nil { return err }
        NewPeer(conn, m)
    }
}
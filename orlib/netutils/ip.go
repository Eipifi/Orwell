package netutils
import "net"

type IpFinder interface {
    Find() (net.IP, error)
}
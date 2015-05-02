package netutils
import (
    "net"
)

var finders []IpFinder = []IpFinder{
    &HttpRawIpFinder{"https://wtfismyip.com/text"},
    &HttpRawIpFinder{"http://myexternalip.com/raw"},
}

func FindExternalIp() net.IP {
    for _, f := range finders {
        ip, _ := f.Find()
        if ip != nil { return ip }
    }
    return nil
}
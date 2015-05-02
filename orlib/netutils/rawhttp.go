package netutils
import (
    "net"
    "net/http"
    "io/ioutil"
    "errors"
    "strings"
)

var ErrParseIpFailed = errors.New("Failed to parse IP response")

type HttpRawIpFinder struct {
    address string
}

func (f *HttpRawIpFinder) Find() (ip net.IP, err error) {
    response, err := http.Get(f.address)
    if err != nil { return }
    raw, err := ioutil.ReadAll(response.Body)
    if err != nil { return }
    rawip := strings.TrimSpace(string(raw))
    ip = net.ParseIP(rawip)
    if ip == nil { err = ErrParseIpFailed }
    return
}

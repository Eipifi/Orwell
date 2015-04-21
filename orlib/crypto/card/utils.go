package card
import (
    "encoding/asn1"
    "errors"
    "io"
    "io/ioutil"
)

func UnmarshalAll(b []byte, p interface{}) error {
    r, e := asn1.Unmarshal(b, p)
    if e != nil { return e}
    if len(r) > 0 { return errors.New("Unparsed bytes remaining") }
    return nil
}

func UnmarshalReader(r io.Reader) (*Card, error) {
    b, err := ioutil.ReadAll(r)
    if err != nil { return nil, err }
    c := &Card{}
    return c, c.UnmarshalPEM(b)
}
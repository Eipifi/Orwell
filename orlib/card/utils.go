package card
import (
    "io"
    "io/ioutil"
    "encoding/pem"
    "errors"
)

func FromPEM(raw []byte) (*Card, error) {
    b, _ := pem.Decode(raw)
    if b == nil {
        return nil, errors.New("Failed to read the PEM file format")
    }
    return Unmarshal(b.Bytes)
}

func FromReader(r io.Reader) (*Card, error) {
    data, err := ioutil.ReadAll(r)
    if err != nil { return nil, err }
    return FromPEM(data)
}

func FromFile(path string) (*Card, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil { return nil, err }
    return FromPEM(data)
}
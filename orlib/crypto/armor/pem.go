package armor
import (
    "encoding"
    "encoding/pem"
    "errors"
)

func Marshal(name string, d encoding.BinaryMarshaler) (b []byte, err error) {
    block := &pem.Block{}
    if block.Bytes, err = d.MarshalBinary(); err != nil { return }
    block.Type = name
    return pem.EncodeToMemory(block), nil
}

func Unmarshal(b []byte, d encoding.BinaryUnmarshaler) (err error) {
    block, r := pem.Decode(b)
    if block == nil { return errors.New("Failed to parse PEM format") }
    if len(r) > 0 { return errors.New("Unparsed bytes remaining") }
    return d.UnmarshalBinary(block.Bytes)
}
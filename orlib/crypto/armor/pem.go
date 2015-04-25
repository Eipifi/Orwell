package armor
import (
    "encoding/pem"
    "errors"
    "orwell/orlib/butils"
)

func Encode(w butils.ByteWritable, name string) ([]byte, err error) {
    block := &pem.Block{}
    if block.Bytes, err = w.WriteBytes(); err != nil { return }
    block.Type = name
    return pem.EncodeToMemory(block), nil
}

func Decode(r butils.ByteReadable, data []byte) error {
    block, rest := pem.Decode(data)
    if block == nil { return errors.New("PEM decode failed") }
    if len(rest) > 0 { return errors.New("Bytes remaining") }
    return r.ReadBytes(block.Bytes)
}
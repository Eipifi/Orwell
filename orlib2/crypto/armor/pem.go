package armor
import (
    "orwell/orlib2"
    "encoding/pem"
    "errors"
)

func Encode(w orlib2.ByteWritable, name string) ([]byte, err error) {
    block := &pem.Block{}
    if block.Bytes, err = w.WriteBytes(); err != nil { return }
    block.Type = name
    return pem.EncodeToMemory(block), nil
}

func Decode(r orlib2.ByteReadable, data []byte) error {
    block, rest := pem.Decode(data)
    if block == nil { return errors.New("PEM decode failed") }
    if len(rest) > 0 { return errors.New("Bytes remaining") }
    return r.ReadBytes(block.Bytes)
}
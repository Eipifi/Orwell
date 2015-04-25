package armor
import (
    "encoding/pem"
    "errors"
    "orwell/orlib/butils"
    "io"
    "io/ioutil"
)

const (
    TextCard = "ORWELL CARD"
    TextKey = "ORWELL PRIVATE KEY"
)

func Encode(data []byte, name string) []byte {
    block := &pem.Block{}
    block.Type = name
    block.Bytes = data
    return pem.EncodeToMemory(block)
}

func Decode(data []byte) ([]byte, error) {
    block, rest := pem.Decode(data)
    if block == nil { return errors.New("PEM decode failed") }
    if len(rest) > 0 { return errors.New("Bytes remaining") }
    return block.Bytes, nil
}

func EncodeObjTo(r butils.ByteWritable, name string, target io.Writer) error {
    buf, err := butils.WriteToBytes(r)
    if err != nil { return err }
    return EncodeBytesTo(buf, name, target)
}

func EncodeBytesTo(data []byte, name string, target io.Writer) error {
    buf := Encode(data, name)
    return butils.WriteFull(target, buf)
}

func DecodeTo(data []byte, r butils.ByteReadable) error {
    buf, err := Decode(data)
    if err != nil { return err }
    return r.ReadBytes(buf)
}

func DecodeAll(r io.Reader) ([]byte, error) {
    data, err := ioutil.ReadAll(r)
    if err != nil { return }
    return Decode(data)
}
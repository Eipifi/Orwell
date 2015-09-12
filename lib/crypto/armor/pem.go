package armor
import (
    "encoding/pem"
    "io/ioutil"
    "errors"
    "orwell/lib/butils"
)

func Read(r butils.ByteReadable, data []byte) error {
    block, _ := pem.Decode(data)
    if block == nil { return errors.New("Failed to parse PEM block") }
    return r.ReadBytes(block.Bytes)
}

func Write(w butils.ByteWritable, title string) (res []byte, err error) {
    b := pem.Block{}
    b.Type = title
    b.Bytes, err = w.WriteBytes()
    if err != nil { return }
    return pem.EncodeToMemory(&b), nil
}

func ReadFromFile(r butils.ByteReadable, path string) error {
    file, err := ioutil.ReadFile(path)
    if err != nil { return err }
    return Read(r, file)
}

func WriteToFile(w butils.ByteWritable, title string, path string) error {
    res, err := Write(w, title)
    if err != nil { return err }
    return ioutil.WriteFile(path, res, 0700)
}
package butils
import (
    "bytes"
    "errors"
)

func ReadAllInto(r Readable, data []byte) (err error) {
    buf := bytes.NewBuffer(data)
    if err = r.Read(buf); err != nil { return }
    if len(buf.Bytes()) > 0 { return errors.New("Unread bytes remaining") }
    return
}

func WriteToBytes(w Writable) ([]byte, error) {
    buf := &bytes.Buffer{}
    if err := w.Write(buf); err != nil { return nil, err }
    return buf.Bytes(), nil
}
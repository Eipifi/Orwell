package butils
import (
    "bytes"
    "errors"
    "io"
    "reflect"
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

func WriteOptional(w io.Writer, target Writable) (err error) {
    if reflect.ValueOf(target).IsNil() { // ugh...
        return WriteByte(w, 0x00)
    } else {
        if err = WriteByte(w, 0x01); err != nil { return }
        return target.Write(w)
    }
}

func ReadOptional(r io.Reader, target Readable) (flag byte, err error) {
    if flag, err = ReadByte(r); err != nil { return }
    if flag == 0x00 { return }
    if flag != 0x01 { return flag, errors.New("Invalid flag value") }
    return flag, target.Read(r)
}
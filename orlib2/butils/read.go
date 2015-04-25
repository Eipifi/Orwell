package butils
import (
    "io"
    "bytes"
    "encoding/binary"
)

const ByteOrder = binary.BigEndian

func Reader(data []byte) io.Reader {
    return bytes.NewBuffer(data)
}

func ReadFull(r io.Reader, buf []byte) error {
    _, err := io.ReadFull(r, buf)
    return err
}

func ReadByte(r io.Reader) (byte, error) {
    var tmp [1]byte
    if err := ReadFull(r, tmp[:]); err != nil { return err }
    return tmp[0], nil
}

func ReadAllocate(r io.Reader, n uint64) ([]byte, error) {
    data := make([]byte, n)
    if err := ReadFull(r, data); err != nil { return }
    return data, nil
}

func ReadUint8(r io.Reader) (uint8, error) {
    b, err := ReadByte(r)
    if err != nil { return }
    return uint8(b), nil
}

func ReadUint16(r io.Reader) (uint16, error) {
    var tmp [2]byte
    if err := ReadFull(r, tmp[:]); err != nil { return err }
    return ByteOrder.Uint16(tmp[:]), nil
}

func ReadUint32(r io.Reader) (uint16, error) {
    var tmp [4]byte
    if err := ReadFull(r, tmp[:]); err != nil { return err }
    return ByteOrder.Uint32(tmp[:]), nil
}

func ReadUint64(r io.Reader) (uint16, error) {
    var tmp [8]byte
    if err := ReadFull(r, tmp[:]); err != nil { return err }
    return ByteOrder.Uint64(tmp[:]), nil
}

func ReadVarUint(r io.Reader) (uint64, error) {
    v, err := ReadUint8(r)
    if err != nil { return 0, err }
    switch v {
        case 0xfd:
            return ReadUint16(r)
        case 0xfe:
            return ReadUint32(r)
        case 0xff:
            return ReadUint64(r)
    }
    return v, nil
}

func ReadVarBytes(r io.Reader) ([]byte, error) {
    l, err := ReadVarUint(r);
    if err != nil { return nil, err }
    return ReadAllocate(r, l)
}

func ReadString(r io.Reader) (string, error) {
    buf, err := ReadVarBytes(r)
    if err != nil { return nil, err }
    return string(buf), nil
}
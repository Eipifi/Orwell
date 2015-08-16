package butils
import (
    "io"
    "bytes"
    "encoding/binary"
    "errors"
)

var ByteOrder = binary.BigEndian

func Reader(data []byte) io.Reader {
    return bytes.NewBuffer(data)
}

func ReadFull(r io.Reader, buf []byte) error {
    _, err := io.ReadFull(r, buf)
    return err
}

func ReadByte(r io.Reader) (byte, error) {
    var tmp [1]byte
    if err := ReadFull(r, tmp[:]); err != nil { return 0, err }
    return tmp[0], nil
}

func ReadAllocate(r io.Reader, n uint64) (data []byte, err error) {
    data = make([]byte, n)
    err = ReadFull(r, data)
    return
}

func ReadUint8(r io.Reader) (uint8, error) {
    b, err := ReadByte(r)
    return uint8(b), err
}

func ReadUint16(r io.Reader) (uint16, error) {
    var tmp [2]byte
    if err := ReadFull(r, tmp[:]); err != nil { return 0, err }
    return ByteOrder.Uint16(tmp[:]), nil
}

func ReadUint32(r io.Reader) (uint32, error) {
    var tmp [4]byte
    if err := ReadFull(r, tmp[:]); err != nil { return 0, err }
    return ByteOrder.Uint32(tmp[:]), nil
}

func ReadUint64(r io.Reader) (uint64, error) {
    var tmp [8]byte
    if err := ReadFull(r, tmp[:]); err != nil { return 0, err }
    return ByteOrder.Uint64(tmp[:]), nil
}

func BytesToUint64(data []byte) uint64 {
    if len(data) != 8 { panic("Invalid slice length for uint64") }
    return ByteOrder.Uint64(data)
}

func ReadVarUint(r io.Reader) (uint64, error) {
    // Integer bounds are enforced in order to ensure that each number has exactly one representation.
    // Even though technically correct, this could create a problem with different block hashes for different integer representations.
    v, err := ReadUint8(r)
    if err != nil { return 0, err }
    switch v {
        case 0xfd:
            v, err := ReadUint16(r)
            if err != nil { return 0, err }
            if v < 0xfd { return 0, errors.New("Illegal VarUint format") }
            return uint64(v), err
        case 0xfe:
            v, err := ReadUint32(r)
            if err != nil { return 0, err }
            if v <= 0xffff { return 0, errors.New("Illegal VarUint format") }
            return uint64(v), err
        case 0xff:
            v, err := ReadUint64(r)
            if err != nil { return 0, err }
            if v <= 0xffffffff { return 0, errors.New("Illegal VarUint format") }
            return v, err
    }
    return uint64(v), nil
}

func ReadVarBytes(r io.Reader, limit uint64) ([]byte, error) {
    l, err := ReadVarUint(r);
    if err != nil { return nil, err }
    if l >= limit { return nil, ErrLimitExceeded }
    return ReadAllocate(r, l)
}

func ReadString(r io.Reader, limit uint64) (string, error) {
    buf, err := ReadVarBytes(r, limit)
    if err != nil { return "", err }
    return string(buf), nil
}
package butils
import "io"

func WriteFull(w io.Writer, buf []byte) error {
    _, err := w.Write(buf)
    return err
}

func WriteByte(w io.Writer, b byte) error {
    return WriteFull(w, []byte{b})
}

func WriteUint8(w io.Writer, val uint8) error {
    return WriteByte(w, byte(val))
}

func WriteUint16(w io.Writer, val uint16) error {
    var tmp [2]byte
    ByteOrder.PutUint16(tmp[:], val)
    return WriteFull(w, tmp[:])
}

func WriteUint32(w io.Writer, val uint32) error {
    var tmp [4]byte
    ByteOrder.PutUint32(tmp[:], val)
    return WriteFull(w, tmp[:])
}

func WriteUint64(w io.Writer, val uint64) error {
    var tmp [8]byte
    ByteOrder.PutUint64(tmp[:], val)
    return WriteFull(w, tmp[:])
}

func Uint64ToBytes(val uint64) []byte {
    var tmp [8]byte
    ByteOrder.PutUint64(tmp[:], val)
    return tmp[:]
}

func WriteVarUint(w io.Writer, val uint64) (err error) {
    if val < 0xfd {
        return WriteUint8(w, uint8(val))
    }

    if val <= 0xffff {
        if err = WriteUint8(w, 0xfd); err != nil { return }
        return WriteUint16(w, uint16(val))
    }

    if val <= 0xffffffff {
        if err = WriteUint8(w, 0xfe); err != nil { return }
        return WriteUint32(w, uint32(val))
    }

    if err = WriteUint8(w, 0xff); err != nil { return }
    return WriteUint64(w, uint64(val))
}

func WriteVarBytes(w io.Writer, buf []byte, limit uint64) error {
    if uint64(len(buf)) >= limit { return ErrLimitExceeded }
    if err := WriteVarUint(w, uint64(len(buf))); err != nil { return err }
    return WriteFull(w, buf)
}

func WriteString(w io.Writer, val string, limit uint64) error {
    return WriteVarBytes(w, []byte(val), limit)
}
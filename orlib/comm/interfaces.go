package comm
import "encoding/binary"

var ByteOrder = binary.BigEndian

type Writable interface {
    Write(w *Writer)
}

type Readable interface {
    Read(r *Reader) error
}

type Msg interface {
    Readable
    Writable
}
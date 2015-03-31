package protocol
import (
    "io"
    "bytes"
)


type Writer struct {
    b bytes.Buffer
    tmp [8]byte
}

func NewWriter() *Writer {
    return &Writer{}
}

func (w *Writer) Write(data []byte) {
    if len(data) == 0 { return } // TODO: find out if optional
    w.b.Write(data)
}

func (w *Writer) Commit(r io.Writer) error {
    _, err := w.b.WriteTo(r)
    return err
}

func (w *Writer) Peek() []byte {
    return w.b.Bytes()
}

func (w *Writer) writeTmp(n int) {
    w.Write(w.tmp[:n])
}

func (w *Writer) WriteUint8(v uint8) {
    w.tmp[0] = byte(v)
    w.writeTmp(1)
}

func (w *Writer) WriteUint16(v uint16) {
    ByteOrder.PutUint16(w.tmp[:], v)
    w.writeTmp(2)
}

func (w *Writer) WriteUint32(v uint32) {
    ByteOrder.PutUint32(w.tmp[:], v)
    w.writeTmp(4)
}

func (w *Writer) WriteUint64(v uint64) {
    ByteOrder.PutUint64(w.tmp[:], v)
    w.writeTmp(8)
}

func (w *Writer) WriteVaruint(v uint64) {
    if v < 0xfd {
        w.WriteUint8(uint8(v))
        return
    }

    if v <= 0xffff {
        w.WriteUint8(0xfd)
        w.WriteUint16(uint16(v))
        return
    }

    if v <= 0xffffffff {
        w.WriteUint8(0xfe)
        w.WriteUint32(uint32(v))
        return
    }

    w.WriteUint8(0xff)
    w.WriteUint64(v)
}

func (w *Writer) WriteVarBytes(data []byte) {
    w.WriteVaruint(uint64(len(data)))
    w.Write(data)
}

func (w *Writer) WriteString(s string) {
    w.WriteVarBytes([]byte(s))
}

func (w *Writer) WriteAddress(a *Address) {
    w.Write(a.IP[:])
    w.WriteUint16(a.Port)
    w.WriteUint64(a.Nonce)
}
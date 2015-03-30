package protocol
import (
    "io"
    "encoding/binary"
    "errors"
    "bytes"
    "orwell/orlib/sig"
)

var ByteOrder = binary.BigEndian
const MaxEnvelopeLength = 65535

type Reader struct {
    r io.Reader
    tmp [8]byte
}

func NewReader(r io.Reader) *Reader {
    return &Reader{r: r}
}

func NewBytesReader(data []byte) *Reader {
    return &Reader{r: bytes.NewReader(data)}
}

func (r *Reader) ReadTo(target []byte) error {
    _, err := io.ReadFull(r.r, target)
    return err
}

func (r *Reader) readToTmp(n int) error {
    return r.ReadTo(r.tmp[:n])
}

func (r *Reader) ReadAllocate(n uint64) ([]byte, error) {
    buf := make([]byte, n)
    err := r.ReadTo(buf)
    if err != nil { return nil, err }
    return buf, nil
}

func (r *Reader) ReadUint8() (uint8, error) {
    if err := r.readToTmp(1); err != nil { return 0, err }
    return uint8(r.tmp[0]), nil
}

func (r *Reader) ReadUint16() (uint16, error) {
    if err := r.readToTmp(2); err != nil { return 0, err }
    return ByteOrder.Uint16(r.tmp[:2]), nil
}

func (r *Reader) ReadUint32() (uint32, error) {
    if err := r.readToTmp(4); err != nil { return 0, err }
    return ByteOrder.Uint32(r.tmp[:4]), nil
}

func (r *Reader) ReadUint64() (uint64, error) {
    if err := r.readToTmp(8); err != nil { return 0, err }
    return ByteOrder.Uint64(r.tmp[:8]), nil
}

func (r *Reader) ReadVaruint() (val uint64, err error) {
    if err = r.ReadTo(r.tmp[:1]); err != nil { return }
    switch r.tmp[0] {
        case 0xfd:
            var t uint16
            t, err = r.ReadUint16()
            val = uint64(t)
        case 0xfe:
            var t uint32
            t, err = r.ReadUint32()
            val = uint64(t)
        case 0xff:
            val, err = r.ReadUint64()
        default:
            val = uint64(r.tmp[0])
    }
    return
}

func (r *Reader) ReadVarBytes() ([]byte, error) {
    l, err := r.ReadVaruint()
    if err != nil { return nil, err }
    return r.ReadAllocate(l)
}

func (r *Reader) ReadStr() (string, error) {
    b, err := r.ReadVarBytes()
    if err != nil { return "", err }
    return string(b), nil
}

func (r *Reader) ReadAddress() (addr *Address, err error) {
    addr = &Address{}
    if err = r.ReadTo(addr.IP[:]); err != nil { return }
    if addr.Port, err = r.ReadUint16(); err != nil { return }
    if addr.Nonce, err = r.ReadUint64(); err != nil { return }
    return
}

func (r *Reader) ReadID() (id *sig.ID, err error) {
    id = &sig.ID
    _, err = r.ReadTo(id[:])
    return
}

func (r *Reader) ReadFrame() (frame *Frame, err error) {
    frame = &Frame{}
    if frame.Command, err = r.ReadVaruint(); err != nil { return }
    if frame.Payload, err = r.ReadVarBytes(); err != nil { return }
    return
}
package db
import (
    "io"
    "orwell/lib/butils"
)

const max_buf_size uint64 = 1 << 32

type RollbackOp struct {
    Bucket []byte
    Key []byte
    Value []byte
}

func (e *RollbackOp) Read(r io.Reader) (err error) {
    if e.Bucket, err = butils.ReadVarBytes(r, max_buf_size); err != nil { return }
    if e.Key, err = butils.ReadVarBytes(r, max_buf_size); err != nil { return }
    flag, err := butils.ReadByte(r)
    if err != nil { return }
    if flag != 0x00 {
        if e.Value, err = butils.ReadVarBytes(r, max_buf_size); err != nil { return }
    }
    return
}

func (e *RollbackOp) Write(w io.Writer) (err error) {
    if err = butils.WriteVarBytes(w, e.Bucket, max_buf_size); err != nil { return }
    if err = butils.WriteVarBytes(w, e.Key, max_buf_size); err != nil { return }
    if e.Value == nil {
        return butils.WriteByte(w, 0x00)
    }
    if err = butils.WriteByte(w, 0x01); err != nil { return }
    if err = butils.WriteVarBytes(w, e.Value, max_buf_size); err != nil { return }
    return
}
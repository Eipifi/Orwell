package orcache
import (
    "io"
    "orwell/lib/butils"
)

const SIZE_KEY uint64 = 32
const SIZE_TYPE uint64 = 32
const SIZE_VALUE uint64 = 4096

type Entry struct {
    Key []byte
    Type []byte
    Value []byte
}

func (e *Entry) Read(r io.Reader) (err error) {
    if e.Key, err = butils.ReadVarBytes(r, SIZE_KEY); err != nil { return }
    if e.Type, err = butils.ReadVarBytes(r, SIZE_TYPE); err != nil { return }
    if e.Value, err = butils.ReadVarBytes(r, SIZE_VALUE); err != nil { return }
    return
}

func (e *Entry) Write(w io.Writer) (err error) {
    if err = butils.WriteVarBytes(w, e.Key, SIZE_KEY); err != nil { return }
    if err = butils.WriteVarBytes(w, e.Type, SIZE_TYPE); err != nil { return }
    if err = butils.WriteVarBytes(w, e.Value, SIZE_VALUE); err != nil { return }
    return
}

func NewEntry(key, tpe, value string) Entry {
    return Entry{Key: []byte(key), Type: []byte(tpe), Value: []byte(value)}
}
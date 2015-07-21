package butils

import (
    "io"
    "errors"
)

/*
    Main byte conversion interfaces for Orwell.
    Types implementing these interfaces can infer the byte length during the read process.
    This means that the object size is either fixed or deducible (eg. by framing protocol, stop character).

    Implementing types shall not check if there are remaining bytes after the read.
    This is the user's responsibility.
*/

type Readable interface {
    Read(io.Reader) error
}

type Writable interface {
    Write(io.Writer) error
}

type Chunk interface {
    Readable
    Writable
}

type ByteReadable interface {
    ReadBytes([]byte) error
}

type ByteWritable interface {
    WriteBytes() ([]byte, error)
}

type ByteChunk interface {
    ByteReadable
    ByteWritable
}

//////////////////////////////

var ErrLimitExceeded = errors.New("Length limit exceeded")
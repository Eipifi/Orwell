package orchain
import (
    "orwell/lib/foo"
    "io"
    "orwell/lib/butils"
    "errors"
)

type MsgGetBlock struct {
    ID foo.U256
}

func (m *MsgGetBlock) Read(r io.Reader) error {
    return m.ID.Read(r)
}

func (m *MsgGetBlock) Write(w io.Writer) error {
    return m.ID.Write(w)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type MsgBlock struct {
    Block *Block
}

func (m *MsgBlock) Read(r io.Reader) error {
    flag, err := butils.ReadByte(r)
    if err != nil { return err }
    if flag == 0x00 {
        m.Block = nil
        return nil
    }
    if flag == 0x01 {
        m.Block = &Block{}
        return m.Block.Read(r)
    }
    return errors.New("Unknown flag value")
}

func (m *MsgBlock) Write(w io.Writer) error {
    if m.Block == nil {
        return butils.WriteByte(w, 0x00)
    } else {
        if err := butils.WriteByte(w, 0x01); err != nil { return err }
        return m.Block.Write(w)
    }
}
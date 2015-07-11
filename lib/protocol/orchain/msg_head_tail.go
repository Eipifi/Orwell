package orchain
import (
    "orwell/lib/butils"
    "io"
)

const MSG_CODE_HEAD byte = 0x01
const MSG_CODE_TAIL byte = 0x81

type MsgHead struct {
    Work butils.Uint256
    Id butils.Uint256
}

func (m *MsgHead) Read(r io.Reader) (err error) {
    if err = m.Work.Read(r); err != nil { return }
    if err = m.Id.Read(r); err != nil { return }
    return nil
}

func (m *MsgHead) Write(w io.Writer) (err error) {
    if err = m.Work.Write(w); err != nil { return }
    if err = m.Id.Write(w); err != nil { return }
    return nil
}

//////////////////////////////////////////////////////

type MsgTail struct {
    Work butils.Uint256
    Headers []Header
}

func (m *MsgTail) Read(r io.Reader) (err error) {
    if err = m.Work.Read(r); err != nil { return }
    var num uint64
    if num, err = butils.ReadVarUint(r); err != nil { return }
    m.Headers = make([]Header, num)
    for i := 0; i < int(num); i += 1 {
        if err = m.Headers[i].Read(r); err != nil { return }
    }
    return nil
}

func (m *MsgTail) Write(w io.Writer) (err error) {
    if err = m.Work.Write(w); err != nil { return }
    if err = butils.WriteVarUint(w, uint64(len(m.Headers))); err != nil { return }
    for _, h := range m.Headers {
        if err = h.Write(w); err != nil { return }
    }
    return
}
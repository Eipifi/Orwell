package comm
import "net"

type MsgReadFunc func(*Reader, Msg) error
type MsgAllocFunc func(*Reader) (Msg, error)
type MsgWriteFunc func(*Writer, Msg)

type Messenger struct {
    conn net.Conn
    r *Reader
    w *Writer
    rf MsgReadFunc
    wf MsgWriteFunc
    af MsgAllocFunc
}

func NewMessager(conn net.Conn, rf MsgReadFunc, wf MsgWriteFunc, af MsgAllocFunc) *Messenger {
    return &Messenger{conn: conn, r: NewReader(conn), w: NewWriter(), rf: rf, wf: wf, af: af}
}

func (ms *Messenger) Read(m Msg) error {
    return ms.rf(ms.r, m)
}

func (ms *Messenger) ReadAny() (Msg, error) {
    return ms.af(ms.r)
}

func (ms *Messenger) Write(m Msg) error {
    ms.wf(ms.w, m)
    return ms.w.Commit(ms.conn)
}

func (ms *Messenger) WriteMany(msgs ...Msg) error {
    for _, m := range msgs {
        ms.wf(ms.w, m)
    }
    return ms.w.Commit(ms.conn)
}

func (ms *Messenger) Close() error {
    return ms.conn.Close()
}
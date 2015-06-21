package obp
import (
    "net"
    "sync/atomic"
    "errors"
    "sync"
    "log"
    "orwell/lib/logging"
)

type Handler func([]byte) ([]byte, error)
var ErrSocketClosed = errors.New("Socket closed")

type Conn struct {
    log *log.Logger
    socket net.Conn
    closed int32
    context uint64
    upstream chan *Frame
    dnstream chan *Frame
    mtx *sync.Mutex
    queries map[uint64] chan *Frame
}

func New(socket net.Conn) *Conn {
    prefix := socket.RemoteAddr().String() + " "
    c := &Conn{}
    c.log = logging.GetLogger(prefix)
    c.socket = socket
    c.closed = 0
    c.context = 0
    c.upstream = make(chan *Frame)
    c.dnstream = make(chan *Frame)
    c.mtx = &sync.Mutex{}
    c.queries = make(map[uint64] chan *Frame)
    go c.sender()
    go c.receiver()
    return c
}

func (c *Conn) sender() {
    defer c.Close()
    for {
        f, ok := <- c.upstream
        if !ok { return }
        err := f.Write(c.socket)
        if err != nil {
            c.log.Println("Failed to write frame: %s", err)
            return
        }
        c.log.Println("Sent %v", f)
    }
}

func (c *Conn) receiver() {
    defer c.Close()
    for {
        f := &Frame{}
        err := f.Read(c.socket)
        if err != nil {
            c.log.Println("Failed to read frame: %s", err)
            return
        }
        c.log.Println("Received %v", f)
        if f.Context % 2 == 0 {
            c.dnstream <- f // may fail
        } else {
            c.mtx.Lock()
            if rc, ok := c.queries[f.Context]; ok {
                delete(c.queries, f.Context)
                c.mtx.Unlock()
                rc <- f // may fail
            } else {
                c.mtx.Unlock()
                c.log.Println("Received frame does not match any asked question")
                return
            }
        }
    }
}

func (c *Conn) Close() {
    if atomic.CompareAndSwapInt32(&(c.closed), 0, 1) {
        c.socket.Close()
        close(c.upstream)
        close(c.dnstream)
        c.mtx.Lock()
        defer c.mtx.Unlock()
        for _, rc := range c.queries {
            close(rc)
        }
        c.log.Println("Closed connection")
    }
}

func (c *Conn) Query(request []byte) (response []byte, err error) {
    f := &Frame{}
    f.Context = atomic.AddUint64(&c.context, 2)
    f.Payload = request
    rc := make(chan *Frame)
    c.mtx.Lock()
    c.queries[f.Context] = rc
    c.mtx.Unlock()
    if response, ok := <- rc; ok {
        return response.Payload, nil
    } else {
        return nil, ErrSocketClosed
    }
}

func (c *Conn) Handle(handler Handler) error {
    if f, ok := <- c.dnstream; ok {
        if response, err := handler(f.Payload); err == nil {
            f.Context += 1
            f.Payload = response
            c.upstream <- f // may fail
            return nil
        } else {
            c.Close()
            return err
        }
    } else { return ErrSocketClosed }
}

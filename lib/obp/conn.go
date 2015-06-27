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
    defer c.log.Println("Closing the sender()")
    for {
        f, ok := <- c.upstream
        if ! ok { return }
        err := f.Write(c.socket)
        if err != nil {
            c.log.Printf("Failed to write frame: %s", err)
            return
        }
        c.log.Println("Sent %v", f)
    }
}

func (c *Conn) receiver() {
    defer c.Close()
    defer c.log.Println("Closing the receiver()")
    for {
        f := &Frame{}
        err := f.Read(c.socket)
        if err != nil {
            c.log.Printf("Failed to read frame: %s", err)
            return
        }
        c.log.Printf("Received %+v", f)
        if f.Context % 2 == 0 {
            maybeWrite(c.dnstream, f)
        } else {
            c.mtx.Lock()
            if rc, ok := c.queries[f.Context-1]; ok {
                delete(c.queries, f.Context-1)
                c.mtx.Unlock()
                rc <- f
            } else {
                c.mtx.Unlock()
                c.log.Println("Received response frame does not match any sent request")
                return
            }
        }
    }
}

func (c *Conn) Close() {
    if atomic.CompareAndSwapInt32(&(c.closed), 0, 1) {
        c.socket.Close()
        close(c.dnstream)
        close(c.upstream)
        c.mtx.Lock()
        defer c.mtx.Unlock()
        for _, rc := range c.queries {
            rc <- nil
        }
        c.log.Println("Closed the connection")
    }
}

func (c *Conn) Query(request []byte) ([]byte, error) {
    f := &Frame{}
    f.Context = atomic.AddUint64(&c.context, 2)
    f.Payload = request
    rc := make(chan *Frame)
    c.mtx.Lock()
    c.queries[f.Context] = rc
    c.mtx.Unlock()
    maybeWrite(c.upstream, f)
    if response := <- rc; response != nil {
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
            maybeWrite(c.upstream, f)
            return nil
        } else {
            c.Close()
            return err
        }
    } else { return ErrSocketClosed }
}

// HACK - we need a way to *maybe* write, and fail silently if the channel is closed
func maybeWrite(c chan *Frame, f *Frame) {
    defer func() { recover() }()
    c <- f
}

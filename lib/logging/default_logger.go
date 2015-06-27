package logging
import (
    "log"
    "os"
    "io"
)

var logTarget WriterProxy = WriterProxy{w: os.Stdout}

type WriterProxy struct {
    w io.Writer
}

func (wp *WriterProxy) Write(p []byte) (n int, err error) {
    return wp.w.Write(p)
}

func GetLogger(prefix string) *log.Logger {
    return log.New(&logTarget, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func DirectTo(w io.Writer) {
    logTarget.w = w
}

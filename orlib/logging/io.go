package logging
import (
    "io"
    "fmt"
)

type PrintReader struct {
    R io.Reader
}

func (d *PrintReader) Read(p []byte) (n int, err error) {
    n, err = d.R.Read(p)
    if err == nil {
        fmt.Println("PrintReader: OK %x\n", p)
    } else {
        fmt.Println("PrintReader: ERR n:%d e:%s\n" , n, err)
    }
    return
}

type PrintWriter struct {
    W io.Writer
}

func (d *PrintWriter) Write(p []byte) (n int, err error) {
    n, err = d.W.Write(p)
    if err == nil {
        fmt.Println("PrintWriter: OK %x\n", p)
    } else {
        fmt.Println("PrintWriter: ERR n:%d e:%s\n" , n, err)
    }
    return
}
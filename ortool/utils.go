package main
import (
    "io"
    "os"
    "io/ioutil"
)

func FileOrSTDIN(path string) (r io.Reader, err error) {
    if path == "" {
        return os.Stdin, nil
    } else {
        return os.Open(path)
    }
}

func FileOrSTDOUT(path string) (r io.Writer, err error) {
    if path == "" {
        return os.Stdout, nil
    } else {
        return os.Open(path)
    }
}

func ReadWholeFileOrSTDIN(path string) (b []byte, err error) {
    var r io.Reader
    if r, err = FileOrSTDIN(path); err != nil { return }
    return ioutil.ReadAll(r)
}

func rs(array []string, idx int) string {
    if idx <= len(array) - 1 {
        return array[idx]
    } else {
        return ""
    }
}
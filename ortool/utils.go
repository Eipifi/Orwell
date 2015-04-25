package main
import (
    "io"
    "os"
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

func rs(array []string, idx int) string {
    if idx <= len(array) - 1 {
        return array[idx]
    } else {
        return ""
    }
}
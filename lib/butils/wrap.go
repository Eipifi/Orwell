package butils

type BRWrapper struct {
    R Readable
}

func (w *BRWrapper) ReadBytes(data []byte) error {
    return ReadAllInto(w.R, data)
}

type BWWrapper struct {
    W Writable
}

func (w *BWWrapper) WriteBytes() ([]byte, error) {
    return WriteToBytes(w.W)
}
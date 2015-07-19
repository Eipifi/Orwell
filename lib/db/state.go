package db
import (
    "orwell/lib/foo"
    "orwell/lib/butils"
    "io"
)

type State struct {
    Length uint64
    Head foo.U256
    Work foo.U256
}

func (s *State) Read(r io.Reader) (err error) {
    if s.Length, err = butils.ReadUint64(r); err != nil { return }
    if err = s.Head.Read(r); err != nil { return }
    if err = s.Work.Read(r); err != nil { return }
    return
}

func (s *State) Write(w io.Writer) (err error) {
    if err = butils.WriteUint64(w, s.Length); err != nil { return }
    if err = s.Head.Write(w); err != nil { return }
    if err = s.Work.Write(w); err != nil { return }
    return
}
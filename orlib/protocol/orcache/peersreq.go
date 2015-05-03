package orcache
import "io"

type PeersReq struct { }

func (*PeersReq) Read(io.Reader) error { return nil }

func (*PeersReq) Write(io.Writer) error { return nil }
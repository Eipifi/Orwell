package orcache
import "io"

type HandshakeAck struct { }

func (m *HandshakeAck) Read(r io.Reader) error { return nil }

func (m *HandshakeAck) Write(w io.Writer) error { return nil }
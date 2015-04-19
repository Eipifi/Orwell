package orcache
import "orwell/orlib/comm"

type HandshakeAck struct { }

func (m *HandshakeAck) Read(r *comm.Reader) error { return nil }

func (m *HandshakeAck) Write(w *comm.Writer) { }
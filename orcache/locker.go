package main
import "orwell/orlib/protocol"

type TokenLocker interface {
    Lock(token protocol.Token) bool
    Unlock(token protocol.Token)
}

type EmptyTokenLocker struct { }

func (l *EmptyTokenLocker) Lock(protocol.Token) bool { return false }

func (l *EmptyTokenLocker) Unlock(protocol.Token) { }
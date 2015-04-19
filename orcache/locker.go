package main
import "orwell/orlib/protocol/types"

type TokenLocker interface {
    Lock(token types.Token) bool
    Unlock(token types.Token)
}

type EmptyTokenLocker struct { }

func (l *EmptyTokenLocker) Lock(types.Token) bool { return false }

func (l *EmptyTokenLocker) Unlock(types.Token) { }
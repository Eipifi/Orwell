package main
import "orwell/orlib/protocol/common"


var Locker TokenLocker = &NullLocker{}

type TokenLocker interface {
    Lock(common.Token) bool
    Unlock(common.Token)
}

type NullLocker struct { }

func (n *NullLocker) Lock(common.Token) bool { return true }

func (n *NullLocker) Unlock(common.Token) { }
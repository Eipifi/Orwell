package main

type TokenLocker interface {
    Lock(token uint64) bool
    Unlock(token uint64)
}
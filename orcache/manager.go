package main
import "orwell/orlib/sig"


type ConnectionManager interface {
    PickPeer(id sig.ID) *Peer
}
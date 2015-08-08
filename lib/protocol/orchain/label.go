package orchain
import "orwell/lib/foo"

/*
    A label can contain:
        - any string
        - ticket
        - announcement
        - transfer
*/

// Hash of the announcement that is to follow
type Ticket foo.U256

// Domain ownership transfer. Proof must be signed by the legal previous owner.
type Transfer struct {
    Domain Domain
    Proof Proof
}

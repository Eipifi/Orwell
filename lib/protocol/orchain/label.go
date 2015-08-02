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

// Information about the new domain owner
type Announcement struct {
    Name string
    Owner foo.U256
    ValidUntil uint64
}

// Domain ownership transfer. Proof must be signed by the legal previous owner.
type Transfer struct {
    Announcement Announcement
    Proof Proof
}

///////////////////////////////

/*
    Block must contain all valid announcements:

    - announcement matching the previously accepted ticket
    - transfer announcement
*/
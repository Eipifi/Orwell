package db
import "orwell/lib/foo"

var BUCKET_TICKET = []byte("ticket")

func (t *Tx) PutTicket(ticket foo.U256) {
    t.Put(BUCKET_TICKET, ticket[:], FLAG)
}

func (t *Tx) IsTicket(ticket foo.U256) bool {
    return t.Get(BUCKET_TICKET, ticket[:]) != nil
}
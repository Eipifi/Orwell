package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/utils"
    "orwell/lib/butils"
    "bytes"
)

var BUCKET_DOMAIN = []byte("domain")
var BUCKET_REGISTERED_DOMAIN = []byte("domain_registered")
var BUCKET_DOMAIN_BLOCK   = []byte("domain_block")

func (t *Tx) GetDomain(id foo.U256) *orchain.Domain {
    d := &orchain.Domain{}
    if t.Read(BUCKET_DOMAIN, id[:], d) {
        return d
    }
    return nil
}

func (t *Tx) PutDomain(d *orchain.Domain) {
    id := d.ID()
    t.Write(BUCKET_DOMAIN, id[:], d)
}

func (t *Tx) GetRegisteredDomain(name string) *orchain.Domain {
    d := &orchain.Domain{}
    if t.Read(BUCKET_REGISTERED_DOMAIN, []byte(name), d) {
        utils.Assert(d.Name == name)
        return d
    }
    return nil
}

func (t *Tx) GetValidRegisteredDomain(name string) (*orchain.Domain) {
    d := t.GetRegisteredDomain(name)
    if d == nil { return nil }
    s := t.GetState()
    if d.ValidUntilBlock >= s.Length {
        return d
    } else {
        return nil
    }
}

func (t *Tx) RegisterDomain(d *orchain.Domain) {
    t.Write(BUCKET_REGISTERED_DOMAIN, []byte(d.Name), d)
}

func (t *Tx) DomainsToRegister(txns []orchain.Transaction) (domains []orchain.Domain) {

    // iterate over transactions, collect transfers and announcements
    just_announced := make(map[foo.U256] orchain.Domain)
    for _, txn := range txns {
        if txn.Payload.Transfer != nil {
            domains = append(domains, txn.Payload.Transfer.Domain)
        }
        if txn.Payload.Domain != nil {
            just_announced[txn.Payload.Domain.ID()] = *txn.Payload.Domain
        }
    }

    var tickets []foo.U256
    // get the tickets announced previously
    state := t.GetState()
    if state.Length >= orchain.BLOCKS_UNTIL_DOMAIN_CONFIRMED {
        prev_bid := t.GetIDByNum(state.Length - orchain.BLOCKS_UNTIL_DOMAIN_CONFIRMED)
        utils.Assert(prev_bid != nil)
        txns := t.GetTransactionsFromBlock(*prev_bid)
        for _, txn := range txns {
            if txn.Payload.Ticket != nil {
                tickets = append(tickets, *txn.Payload.Ticket)
            }
        }
    }

    // for each ticket, check if the domain was announced
    for _, ticket := range tickets {
        domain := t.GetDomain(ticket)
        if domain == nil {
            // check if the domain announcement was included in this particular block, last-minute
            tmp, ok := just_announced[ticket]
            if ok {
                domain = &tmp
            }
        }
        if domain != nil {
            // The domain was announced for the ticket.
            // Now check if the domain can be registered.
            if t.GetValidRegisteredDomain(domain.Name) == nil {
                domains = append(domains, *domain)
            }
        }
    }
    return domains

}

func (t *Tx) IsTransferLegal(transfer *orchain.Transfer) bool {
    if transfer.Proof.CheckWritable(&transfer.Domain) != nil { return false }
    reg_domain := t.GetValidRegisteredDomain(transfer.Domain.Name)
    if reg_domain == nil { return false }
    return reg_domain.Owner == transfer.Proof.PublicKey.ID()
}

func (t *Tx) PutRegisteredDomainsFromBlock(id foo.U256, domains []orchain.Domain) {
    buf := bytes.Buffer{}
    utils.Ensure(butils.WriteSlice(&buf, orchain.BLOCK_DOMAIN_MAX, domains))
    t.Put(BUCKET_DOMAIN_BLOCK, id[:], buf.Bytes())
}

func (t *Tx) GetDomainsFromBlock(id foo.U256) (domains []orchain.Domain) {
    data := t.Get(BUCKET_DOMAIN_BLOCK, id[:])
    if data == nil { return }
    buf := bytes.NewBuffer(data)
    utils.Ensure(butils.ReadSlice(buf, orchain.BLOCK_DOMAIN_MAX, &domains))
    return
}
package db
import (
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "orwell/lib/utils"
)

var BUCKET_DOMAIN = []byte("domain")
var BUCKET_REGISTERED_DOMAIN = []byte("domain_registered")

func (t *Tx) GetDomain(id foo.U256) (d *orchain.Domain) {
    d = &orchain.Domain{}
    t.Read(BUCKET_DOMAIN, id[:], d)
    return
}

func (t *Tx) PutDomain(d *orchain.Domain) {
    id := d.ID()
    t.Write(BUCKET_DOMAIN, id[:], d)
}

func (t *Tx) GetRegisteredDomain(name string) (d *orchain.Domain) {
    d = &orchain.Domain{}
    t.Read(BUCKET_REGISTERED_DOMAIN, []byte(name), d)
    return
}

func (t *Tx) RegisterDomain(d *orchain.Domain) {
    t.Write(BUCKET_REGISTERED_DOMAIN, []byte(d.Name), d)
}

func (t *Tx) DomainsToRegister(txns []orchain.Transaction) (domains []orchain.Domain) {

    // iterate over transactions, collect transfers and announcements
    var just_announced map[foo.U256] orchain.Domain
    for _, txn := range txns {
        if txn.Payload.Transfer != nil {
            domains = append(domains, txn.Payload.Transfer.Domain)
        }
        if txn.Payload.Domain != nil {
            just_announced[txn.Payload.Domain.ID()] = *txn.Payload.Domain
        }
    }

    state := t.GetState()
    var tickets []foo.U256
    // get the tickets announced previously
    if state.Length >= orchain.BLOCKS_UNTIL_DOMAIN_CONFIRMED {
        prev_bid := t.GetIDByNum(state.Length - orchain.BLOCKS_UNTIL_DOMAIN_CONFIRMED)
        utils.Assert(prev_bid != nil)
        prev := t.GetBlock(*prev_bid)
        for _, txn := range prev.Transactions {
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
            if t.IsDomainFree(domain.Name) {
                domains = append(domains, *domain)
            }
        }
    }
    return domains
}

func (t *Tx) IsDomainFree(name string) bool {
    reg_domain := t.GetRegisteredDomain(name)
    if reg_domain != nil {
        utils.Assert(reg_domain.Name == name) // just to be sure
        state := t.GetState()
        if reg_domain.ValidUntilBlock >= state.Length {
            return false
        }
    }
    return true
}

func (t *Tx) IsTransferLegal(transfer *orchain.Transfer) bool {
    if transfer.Proof.CheckObject(&transfer.Domain) != nil { return false }
    reg_domain := t.GetRegisteredDomain(transfer.Domain.Name)
    if reg_domain == nil { return false }
    state := t.GetState()
    if reg_domain.ValidUntilBlock < state.Length { return false }
    return reg_domain.Owner == transfer.Proof.PublicKey.ID()
}
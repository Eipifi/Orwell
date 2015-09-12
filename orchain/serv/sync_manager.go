package serv
import (
    "log"
    "orwell/lib/logging"
    "math/rand"
    "time"
    "orwell/lib/protocol/orchain"
    "orwell/lib/foo"
    "errors"
    "orwell/lib/db"
    "github.com/deckarep/golang-set"
)

const SYNC_MIN_PERIOD = 5
const SYNC_MAX_PERIOD = 10

type SyncManager struct {
    log *log.Logger
    pending_txns mapset.Set
}

func (m *SyncManager) syncLoop() {
    for {
        peers := ConnMgr().GetRandomPeers(1)
        if len(peers) > 0 {
            if err := m.sync(peers[0]); err != nil {
                m.log.Printf("Failed to sync: %v\n", err)
                peers[0].Close()
            }
        }
        interval := SYNC_MIN_PERIOD + rand.Intn(SYNC_MAX_PERIOD - SYNC_MIN_PERIOD)
        time.Sleep(time.Duration(interval) * time.Second)
    }
}

func (m *SyncManager) PushBlock(block *orchain.Block) (err error) {
    db.Get().Update(func(t *db.Tx) {
        err = t.PushBlock(block)
        if err == nil {
            for _, txn := range block.Transactions {
                m.pending_txns.Remove(txn)
            }
        }
    })
    return
}

func (m *SyncManager) syncBlocks(peer *Peer) (err error) {
    rsp := &orchain.MsgTail{}
    var state *db.State
    db.Get().View(func(t *db.Tx) {
        state = t.GetState()
    })
    var revert uint64 = 1
    for len(rsp.Headers) == 0 {
        if rsp, err = peer.AskHead(revert); err != nil { return }
        if foo.Compare(rsp.Work, state.Work) != 1 { return nil }
        if revert > state.Length && len(rsp.Headers) == 0 { return errors.New("Node advertises more work, yet does not send any headers after genesis block") }
        revert *= 2
    }

    headers := rsp.Headers
    db.Get().Update(func(t *db.Tx) {
        // Iterate to first unknown header
        for {
            if len(headers) == 0 { return }
            if t.GetNumByID(headers[0].ID()) == nil { break }
            headers = headers[1:]
        }

        // TODO verify here if headers are properly signed, sum up to the declared work, and make overall sense

        // Drop the obsolete blocks
        for (headers[0].Previous) != (t.GetState().Head) {
            t.PopBlock()
        }
    })

    // Download blocks and apply in order
    for _, h := range headers {
        var block_rsp *orchain.MsgBlock
        if block_rsp, err = peer.AskBlock(h.ID()); err != nil { return }
        if block_rsp.Block == nil { return } // The peer promised to deliver the block, and failed - what to do?
        if err = m.PushBlock(block_rsp.Block); err != nil { return }
    }

    m.log.Printf("Sync successful, downloaded %v blocks", len(headers))
    return
}

func (m *SyncManager) syncTxns(peer *Peer) (error) {
    rsp, err := peer.AskTxns()
    if err != nil { return err }
    db.Get().Update(func(t *db.Tx) {
        for _, txn := range rsp.Transactions {
            if err = t.MaybeStoreUnconfirmedTransaction(&txn); err != nil {
                m.log.Printf("Received txn not stored: %v", err)
            }
            // TODO: maybe log the errors?
        }
    })
    return nil
}

func (m *SyncManager) sync(peer *Peer) (err error) {
    if err = m.syncBlocks(peer); err != nil { return }
    if err = m.syncTxns(peer); err != nil { return }
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var syncInstance *SyncManager

func SyncMgr() *SyncManager { // TODO: synchronize
    if syncInstance == nil {
        syncInstance = &SyncManager{}
        syncInstance.log = logging.GetLogger("")
        syncInstance.pending_txns = mapset.NewSet()
        go syncInstance.syncLoop()
    }
    return syncInstance
}
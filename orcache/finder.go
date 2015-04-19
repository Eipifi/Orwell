package main
import "orwell/orlib/protocol/orcache"


func Find(msg *orcache.Get, env *Env) *GetResponse {

    if data := env.Cache.Get(msg.ID); data != nil {
        return &GetResponse{data, 0}
    }

    if ! env.Locker.Lock(msg.Token) { return &GetResponse{nil, msg.TTL}}
    defer env.Locker.Unlock(msg.Token)

    ttl := msg.TTL
    for {
        if ttl == 0 { break }
        ttl -= 1

        job := &GetJob{}
        job.Msg = &orcache.Get{msg.Token, ttl, msg.ID, msg.Version}
        job.Sink = make(chan *GetResponse, 1)

        peer := env.Manager.PickPeer(job.Msg.ID)
        if peer == nil || !peer.MaybeSendJob(job) { break }

        response := <- job.Sink

        if response == nil {
            // failed to deliver job to a peer
            // for now we end the finder procedure, maybe retry the peer search?
            return &GetResponse{nil, ttl}
        } else {
            if response.Bytes != nil {
                env.Cache.Put(msg.ID, response.Bytes)
                return response
            } else if response.TTL < ttl {
                ttl = response.TTL
            }
            // Todo: recheck the cache
        }
    }
    return &GetResponse{nil, 0}
}
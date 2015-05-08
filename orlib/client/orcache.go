package client
import (
    "io"
    "orwell/orlib/protocol/common"
    "orwell/orlib/protocol/orcache"
    "errors"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/crypto/card"
)

type HsCheck func(*orcache.Handshake) bool
var ErrTokenMismatch = errors.New("Token mismatch")

func ShakeHands(c io.ReadWriter, agent string, address *common.Address, fn HsCheck) (hs *orcache.Handshake, err error) {
    hs = &orcache.Handshake{}
    if err = orcache.WriteMessage(c, &orcache.Handshake{orcache.Magic, orcache.Version, agent, address}); err != nil { return }
    if err = orcache.ReadMessage(c, hs); err != nil { return }
    if hs.Magic != orcache.Magic { return nil, errors.New("Magic value mismatch") }
    if fn != nil && !fn(hs) { return nil, errors.New("Handshake rejected") }
    if err = orcache.WriteMessage(c, &orcache.HandshakeAck{}); err != nil { return }
    if err = orcache.ReadMessage(c, &orcache.HandshakeAck{}); err != nil { return }
    return
}

func Fetch(c io.ReadWriter, id *hash.ID, version uint64) (card *card.Card, err error) {
    token := common.NewRandomToken()
    if err = orcache.WriteMessage(c, &orcache.FetchReq{token, common.MaxTTLValue, id, version}); err != nil { return }
    rsp := &orcache.FetchRsp{}
    if err = orcache.ReadMessage(c, rsp); err != nil { return }
    if rsp.Token != token { return nil, ErrTokenMismatch }
    return rsp.Card, nil
}

func Publish(c io.ReadWriter, card *card.Card) (ttl common.TTL, err error) {
    token := common.NewRandomToken()
    if err = orcache.WriteMessage(c, &orcache.PublishReq{token, common.MaxTTLValue, card}); err != nil { return }
    rsp := &orcache.PublishRsp{}
    if err = orcache.ReadMessage(c, rsp); err != nil { return }
    if rsp.Token != token { return 0, ErrTokenMismatch }
    return rsp.TTL, nil
}
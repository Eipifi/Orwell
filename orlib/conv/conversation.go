package conv
import (
    "orwell/orlib/protocol/orcache"
    "net"
    "orwell/orlib/butils"
    "errors"
    "orwell/orlib/protocol/common"
    "orwell/orlib/crypto/hash"
    "orwell/orlib/crypto/card"
)

const OrcacheMagic = 0xf4eed077
const SupportedVersion = 1

var ErrTokenMismatch = errors.New("Token mismatch")

type Conversation struct {
    Conn net.Conn
    Hs *orcache.Handshake
}

func CreateConversation(conn net.Conn) *Conversation {
    return &Conversation{conn, nil}
}

func CreateTCPConversation(addr string) (*Conversation, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil { return nil, err }
    return CreateConversation(conn), nil
}

func (c *Conversation) Write(chunk butils.Chunk) error {
    msg := orcache.Msg(chunk)
    if msg == nil { return errors.New("Unknown chunk type") }
    return msg.Write(c.Conn)
}

func (c *Conversation) ReadAny() (chunk butils.Chunk, err error) {
    msg := &orcache.Message{}
    if err = msg.Read(c.Conn); err != nil { return }
    return msg.Chunk, nil
}

func (c *Conversation) ReadSpecific(chunk butils.Chunk) (err error) {
    return orcache.Msg(chunk).Read(c.Conn)
}

func (c *Conversation) DoHandshake(userAgent string, addr *common.Address) (err error) {
    c.Hs, err = ShakeHands(c.Conn, userAgent, addr)
    return err
}

func (c *Conversation) DoGet(id *hash.ID, ver uint64) (cd *card.Card, err error) {
    token := common.NewRandomToken()
    if err = c.Write(&orcache.GetReq{token, common.MaxTTLValue, id, ver}); err != nil { return }
    rsp, err := c.ReadAny()
    if err != nil { return nil, err }
    switch r := rsp.(type) {
        case *orcache.GetRsp:
            if r.Token != token { return nil, ErrTokenMismatch }
            return r.Card, nil
    }
    return nil, errors.New("Unexpected response type")
}

func (c *Conversation) DoPut(cd *card.Card) (ttl common.TTL, err error) {
    token := common.NewRandomToken()
    if err = c.Write(&orcache.PublishReq{token, common.MaxTTLValue, cd}); err != nil { return }
    rsp := &orcache.PublishRsp{}
    if err = c.ReadSpecific(rsp); err != nil { return }
    if rsp.Token != token { return common.TTL(0), ErrTokenMismatch }
    return rsp.TTL, nil
}
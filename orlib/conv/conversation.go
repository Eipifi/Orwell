package conv
import (
    "orwell/orlib/protocol/orcache"
    "net"
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

func (c *Conversation) Write(msg orcache.Message) error {
    return orcache.WriteMessage(c.Conn, msg)
}

func (c *Conversation) ReadAny() (msg orcache.Message, err error) {
    return orcache.ReadAnyMessage(c.Conn)
}

func (c *Conversation) ReadSpecific(msg orcache.Message) (err error) {
    return orcache.ReadMessage(c.Conn, msg)
}

func (c *Conversation) DoHandshake(userAgent string, addr *common.Address) (err error) {
    c.Hs, err = ShakeHands(c.Conn, userAgent, addr)
    return err
}

func (c *Conversation) DoGet(id *hash.ID, ver uint64) (cd *card.Card, err error) {
    token := common.NewRandomToken()
    if err = c.Write(&orcache.FetchReq{token, common.MaxTTLValue, id, ver}); err != nil { return }
    rsp, err := c.ReadAny()
    if err != nil { return nil, err }
    switch r := rsp.(type) {
        case *orcache.FetchRsp:
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
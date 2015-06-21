package orchain
import "orwell/lib/obp"

func CodeToMsg(code byte) (obp.Msg, bool) {
    if code == MSG_CODE_HANDSHAKE_REQ { return &HandshakeReq{}, true }
    if code == MSG_CODE_HANDSHAKE_RSP { return &HandshakeRsp{}, true }
    return nil, false
}

func MsgToCode(msg obp.Msg) (byte, bool) {
    switch msg.(type) {
        case (*HandshakeReq): return MSG_CODE_HANDSHAKE_REQ, true
        case (*HandshakeRsp): return MSG_CODE_HANDSHAKE_RSP, true
    }
    return 0x00, false
}
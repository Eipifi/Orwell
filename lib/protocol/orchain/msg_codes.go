package orchain
import "orwell/lib/obp"

func CodeToMsg(code byte) (obp.Msg, bool) {
    if code == MSG_CODE_HEAD { return &MsgHead{}, true }
    if code == MSG_CODE_TAIL { return &MsgTail{}, true }
    if code == MSG_CODE_GET_BLOCK { return &MsgGetBlock{}, true }
    if code == MSG_CODE_BLOCK { return &MsgBlock{}, true }
    if code == MSG_CODE_GET_TXNS { return &MsgGetTxns{}, true }
    if code == MSG_CODE_TXNS { return &MsgTxns{}, true }
    return nil, false
}

func MsgToCode(msg obp.Msg) (byte, bool) {
    switch msg.(type) {
        case (*MsgHead): return MSG_CODE_HEAD, true
        case (*MsgTail): return MSG_CODE_TAIL, true
        case (*MsgGetBlock): return MSG_CODE_GET_BLOCK, true
        case (*MsgBlock): return MSG_CODE_BLOCK, true
        case (*MsgGetTxns): return MSG_CODE_GET_TXNS, true
        case (*MsgTxns): return MSG_CODE_TXNS, true
    }
    return 0x00, false
}

const MSG_CODE_HEAD byte = 0x01
const MSG_CODE_TAIL byte = 0x81
const MSG_CODE_GET_BLOCK byte = 0x02
const MSG_CODE_BLOCK byte = 0x82
const MSG_CODE_GET_TXNS byte = 0x03
const MSG_CODE_TXNS byte = 0x83
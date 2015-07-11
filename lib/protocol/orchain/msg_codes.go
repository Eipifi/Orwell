package orchain
import "orwell/lib/obp"

func CodeToMsg(code byte) (obp.Msg, bool) {
    if code == MSG_CODE_HEAD { return &MsgHead{}, true }
    if code == MSG_CODE_TAIL { return &MsgTail{}, true }
    return nil, false
}

func MsgToCode(msg obp.Msg) (byte, bool) {
    switch msg.(type) {
        case (*MsgHead): return MSG_CODE_HEAD, true
        case (*MsgTail): return MSG_CODE_TAIL, true
    }
    return 0x00, false
}
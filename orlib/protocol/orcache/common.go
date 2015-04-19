package orcache
import (
    "reflect"
    "orwell/orlib/comm"
)

const OrcacheMagic = 0xf4eed077
const SupportedVersion = 1

type msgTypeEntry struct {
    Command uint64
    Type reflect.Type
}

var msgTypes = []msgTypeEntry {
    msgTypeEntry{0x01, reflect.TypeOf(Handshake{})},
    msgTypeEntry{0x02, reflect.TypeOf(Get{})},
    msgTypeEntry{0x81, reflect.TypeOf(HandshakeAck{})},
    msgTypeEntry{0x82, reflect.TypeOf(CardFound{})},
    msgTypeEntry{0x83, reflect.TypeOf(CardNotFound{})},
}

func getMsgCommand(m comm.Msg) uint64 {
    t := reflect.TypeOf(m)
    for _, e := range msgTypes {
        if t == reflect.PtrTo(e.Type) {
            return e.Command
        }
    }
    return 0x00
}

func getCommandMsg(c uint64) comm.Msg {
    for _, e := range msgTypes {
        if e.Command == c {
            return reflect.New(e.Type).Interface().(comm.Msg)
        }
    }
    return nil
}
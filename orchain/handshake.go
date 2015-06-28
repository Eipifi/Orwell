package main
import "orwell/lib/protocol/orchain"

func GenerateHandshake() *orchain.HandshakeReq {
    hs := &orchain.HandshakeReq{}
    hs.Magic = 42
    hs.Fields = make(map[string] string)
    hs.Fields["version"] = "1"
    hs.Fields["head"] = Storage.Head().String()
    return hs
}

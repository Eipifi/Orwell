package orchain
import (
    "io"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
    "orwell/lib/foo"
)

// 112B

type Header struct {
    Previous foo.U256
    MerkleRoot foo.U256
    Difficulty foo.U256
    Timestamp uint64
    Nonce uint64
}

func (h *Header) Read(r io.Reader) (err error) {
    if err = h.Previous.Read(r); err != nil { return }
    if err = h.MerkleRoot.Read(r); err != nil { return }
    if err = h.Difficulty.Read(r); err != nil { return }
    if h.Timestamp, err = butils.ReadUint64(r); err != nil { return }
    if h.Nonce, err = butils.ReadUint64(r); err != nil { return }
    return
}

func (h *Header) Write(w io.Writer) (err error) {
    if err = h.Previous.Write(w); err != nil { return }
    if err = h.MerkleRoot.Write(w); err != nil { return }
    if err = h.Difficulty.Write(w); err != nil { return }
    if err = butils.WriteUint64(w, h.Timestamp); err != nil { return }
    if err = butils.WriteUint64(w, h.Nonce); err != nil { return }
    return
}

func (h *Header) ID() foo.U256 {
    i, _ := hash.HashOf(h) // nothing can go wrong here (...right?)
    return i
}
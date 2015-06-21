package orchain
import (
    "io"
    "orwell/lib/butils"
    "orwell/lib/crypto/hash"
)

// 82B

type Header struct {
    Previous butils.Uint256
    MerkleRoot butils.Uint256
    Timestamp uint64
    Difficulty uint16
    Nonce uint64
}

func (h *Header) Read(r io.Reader) (err error) {
    if err = h.Previous.Read(r); err != nil { return }
    if err = h.MerkleRoot.Read(r); err != nil { return }
    if h.Timestamp, err = butils.ReadUint64(r); err != nil { return }
    if h.Difficulty, err = butils.ReadUint16(r); err != nil { return }
    if h.Nonce, err = butils.ReadUint64(r); err != nil { return }
    return
}

func (h *Header) Write(w io.Writer) (err error) {
    if err = h.Previous.Write(w); err != nil { return }
    if err = h.MerkleRoot.Write(w); err != nil { return }
    if err = butils.WriteUint64(w, h.Timestamp); err != nil { return }
    if err = butils.WriteUint16(w, h.Difficulty); err != nil { return }
    if err = butils.WriteUint64(w, h.Nonce); err != nil { return }
    return
}

func (h *Header) ID() butils.Uint256 {
    i, _ := hash.HashOf(h) // nothing can go wrong here (...right?)
    return i
}
package sig
import (
    "math/big"
    "errors"
    "io"
    "orwell/lib/butils"
)

const SIG_INT_LEN_BYTES = 32
const SIG_INT_LEN_BITS = SIG_INT_LEN_BYTES * 8

type Signature struct {
    R, S *big.Int
}

func (s *Signature) Read(r io.Reader) (err error) {
    if s.R, err = readUint256(r); err != nil { return }
    if s.S, err = readUint256(r); err != nil { return }
    return
}

func (s *Signature) Write(w io.Writer) (err error) {
    if err = writeUint256(w, s.R); err != nil { return }
    if err = writeUint256(w, s.S); err != nil { return }
    return
}

func writeUint256(w io.Writer, b *big.Int) error {
    if b.Sign() < 0 { return errors.New("big.Int must be non-negative") }
    if b.BitLen() > SIG_INT_LEN_BITS { return errors.New("big.Int must be at most 256 bits long") }
    buf := b.Bytes()
    for (len(buf) < SIG_INT_LEN_BYTES) {
        buf = append([]byte{0x00}, buf...) // prepend zero
    }
    return butils.WriteFull(w, buf)
}

func readUint256(r io.Reader) (b *big.Int, err error) {
    buf, err := butils.ReadAllocate(r, uint64(SIG_INT_LEN_BYTES))
    if err != nil { return }
    b = big.NewInt(0)
    b.SetBytes(buf)
    return
}
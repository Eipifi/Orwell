package sig
import (
    "github.com/eipifi/asn1"
    "crypto/x509/pkix"
    "io"
)

type PublicKeyInfo struct {
    Raw       asn1.RawContent
    Algorithm pkix.AlgorithmIdentifier
    PublicKey asn1.BitString
}

func (i *PublicKeyInfo) Read(r io.Reader) error {
    return asn1.UnmarshalFromReader(i, r)
}

func (i *PublicKeyInfo) Write(w io.Writer) error {
    return asn1.MarshalToWriter(i, w)
}
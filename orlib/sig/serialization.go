package sig
import (
    "errors"
    "crypto/x509"
    "crypto/ecdsa"
    "crypto/rsa"
)

type Serializable interface {
    Serialize() []byte
}

func ParsePubKey(derBytes []byte) (PubKey, error) {
    pub, err := x509.ParsePKIXPublicKey(derBytes)
    if err != nil {
        return nil, err
    } else {
        switch key := pub.(type) {
            case *ecdsa.PublicKey:
                return ecdsaPubKey{key}, nil
            case *rsa.PublicKey:
                return nil, errors.New("RSA not yet implemented")
            default:
                return nil, errors.New("sig: only RSA and ECDSA public keys supported")
        }
    }
}

func ParsePrvKey(derBytes []byte) (PrvKey, error) {
    prv, err := x509.ParseECPrivateKey(derBytes)
    if err != nil {
        return nil, err
    } else {
        key := ecdsaPrvKey{prv}
        return key, nil
    }
    return nil, errors.New("Not yet implemented")
}
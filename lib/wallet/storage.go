package wallet
import (
    "io/ioutil"
    "orwell/lib/config"
    "strings"
    "encoding/pem"
    "errors"
)

func Import(file string) (*Wallet, error) {
    w := &Wallet{}
    file_contents, err := ioutil.ReadFile(file)
    if err != nil { return nil, err }
    block, _ := pem.Decode(file_contents)
    if block == nil { return nil, errors.New("Failed to parse PEM block") }
    if err = w.PrvKey.ReadBytes(block.Bytes); err != nil { return nil, err }
    return w, nil
}

func (w *Wallet) Export(file string) error {
    key_contents, err := w.PrvKey.WriteBytes()
    if err != nil { return err }
    pem_contents := pem.EncodeToMemory(&pem.Block{
        Type: "ORWELL PRIVATE KEY",
        Bytes: key_contents,
    })
    return ioutil.WriteFile(file, pem_contents, 0700)
}

func ListWallets() (result []Wallet) {
    files, err := ioutil.ReadDir(config.Path(config.WALLET_DIR))
    if err != nil { return }
    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".pem") {
            w, err := Import(walletPath(f.Name()))
            if err == nil {
                result = append(result, *w)
            }
        }
    }
    return
}

func walletPath(name string) string {
    return config.Path(config.WALLET_DIR + "/" + name)
}

func (w *Wallet) ExportDefault() error {
    return w.Export(walletPath(w.ID().String() + ".pem"))
}
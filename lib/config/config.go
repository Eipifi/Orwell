package config
import (
    "os"
    "log"
    "github.com/mitchellh/go-homedir"
    "stathat.com/c/jconfig"
    "path"
    "orwell/lib/butils"
)

const CONFIG_FILENAME = "config.json"
const WALLET_DIR = "wallets"
const DEFAULT_CONFIG =
`{
    "port": 1984
}`

var cfg struct {
    home string
    vars *jconfig.Config
}

func LoadDefault() {
    path := os.Getenv("ORCHAIN_PATH")
    if path == "" {
        path = "~/.orchain"
    }
    Load(path)
}

func Load(path string) {
    if err := load(path); err != nil {
        log.Fatal(err)
    }
}

func load(p string) (err error) {
    if cfg.home, err = homedir.Expand(p); err != nil { return err }
    if err = os.MkdirAll(cfg.home, 0755); err != nil { return }
    if err = os.MkdirAll(Path(WALLET_DIR), 0700); err != nil { return }
    config_path := Path(CONFIG_FILENAME)
    if _, err := os.Stat(config_path); os.IsNotExist(err) {
        writeDefaultConfig()
    }
    cfg.vars = jconfig.LoadConfig(config_path)
    return nil
}

func Path(p string) string {
    return path.Clean(cfg.home + "/" + p)
}

func Get(key string) string {
    return cfg.vars.GetString(key)
}

func GetInt(key string) int {
    return cfg.vars.GetInt(key)
}

func writeDefaultConfig() error {
    file, err := os.Create(Path(CONFIG_FILENAME))
    if err != nil { return err }
    return butils.WriteFull(file, []byte(DEFAULT_CONFIG))
}
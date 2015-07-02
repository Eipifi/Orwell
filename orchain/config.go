package main
import (
    "os"
    "github.com/mitchellh/go-homedir"
    "stathat.com/c/jconfig"
    "orwell/lib/butils"
)

/*
    - fetch the ORCHAIN_DIR path from env
    - if does not exist, use relative path to user home dir
    - if the ORCHAIN_DIR does not exist, attempt creating it.
    - if the ORCHAIN_DIR does not contain a config file, create it with default contents
    - create database in the the dir
    - create a log file in the dir
    - if all went well, start server and console

*/

const ORCHAIN_DIR_ENV_KEY = "ORCHAIN_DIR"
const ORCHAIN_DIR_DEFAULT = "~/.orchain"
const ORCHAIN_CFG_FILE_NAME = "config.json"

const ORCHAIN_CFG_DEFAULT = "{}"

type ConfigManager struct {
    directory string
    config *jconfig.Config
}

func InitConfig() (cfg *ConfigManager, err error) {
    cfg = &ConfigManager{}
    cfg.directory = os.Getenv(ORCHAIN_DIR_ENV_KEY)

    if cfg.directory == "" {
        cfg.directory = ORCHAIN_DIR_DEFAULT
    }

    // convert the relative directory path to absolute
    if cfg.directory, err = homedir.Expand(cfg.directory); err != nil { return }

    // Create the directory, if necessary
    if err = os.MkdirAll(cfg.directory, 0755); err != nil { return }

    cfgFileName := cfg.RelPath(ORCHAIN_CFG_FILE_NAME)

    // Check if the config file exists
    if _, err = os.Open(cfgFileName); err != nil {
        // File does not exist, create it
        var file *os.File
        if file, err = os.Create(cfgFileName); err != nil { return }
        if err = butils.WriteFull(file, []byte(ORCHAIN_CFG_DEFAULT)); err != nil { return }
    }

    // Load the config file
    cfg.config = jconfig.LoadConfig(cfgFileName)
    return
}

func (c *ConfigManager) RelPath(path string) string {
    return c.directory + "/" + path
}

func (c *ConfigManager) MinerAddress() butils.Uint256 {
    return butils.Uint256{}
}

func (c *ConfigManager) Port() uint16 {
    return 1984
}
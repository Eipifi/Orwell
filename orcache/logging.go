package main
import (
    "log"
    "os"
)

var Info = log.New(os.Stdout, "INFO: ", log.Ldate | log.Ltime | log.Lshortfile)
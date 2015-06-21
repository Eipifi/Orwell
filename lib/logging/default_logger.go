package logging
import (
    "log"
    "os"
)


func GetLogger(prefix string) *log.Logger {
    return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

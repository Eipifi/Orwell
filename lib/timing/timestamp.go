package timestamps
import "time"

func CurrentTimestamp() uint64 {
    return uint64(time.Now().Unix())
}

package together

import (
    "fmt"
    "time"
)

var t_start time.Time
func t_reset() {
    t_start = time.Now()
}

func t_time(args ...interface{}) {
    diff := time.Now().Sub(t_start)
    str  := ""

    for i, arg := range args {
        if i > 0 {
            str += " "
        }
        str += fmt.Sprint(arg)
    }
    
    fmt.Printf("%15s %s\n", diff, str)
}

func t_sleep(c int) {
    for i := 0; i < c; i++ {
        time.Sleep(1)
    }
}
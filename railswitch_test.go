package together

import (
    "fmt"
    "testing"
    "time"
    "sync"
)


func benchmarkRailSwitch(bn, sz int) {
    for i := 0; i < bn; i++ {
        rs := NewRailSwitch()
        wg := sync.WaitGroup{}
        wg.Add(sz)
        for j := 0; j < sz; j++ {
            go func(x int) {
                rs.Queue(x, 1)
                wg.Done()
                _ = x
                rs.Proceed(x)
            }(j)
        }
        wg.Wait()
        rs.Close()
    }
}


func BenchmarkRailSwitch_5(b *testing.B) {
    benchmarkRailSwitch(b.N, 5)
}
func BenchmarkRailSwitch_20(b *testing.B) {
    benchmarkRailSwitch(b.N, 20)
}

var st time.Time
func startTime() {
    st = time.Now()
}

func timed() time.Duration {
    return time.Now().Sub(st)
}

func TestRailSwitch1(t *testing.T) {

    const (
        a = iota
        b
        c
    )

    rs := NewRailSwitch()
    n  := map[int] string {
        a: "A", b: "B", c: "C",
    }
    sleep := func(s int) {
        time.Sleep(time.Millisecond * time.Duration(s))
    }
    do := func(p, s int) {
        sleep(s)
        fmt.Printf("%v %s\n", timed(), n[p])
        rs.Proceed(p)
    }

    startTime()
    // A
    go func() {
        rs.Queue(a, 3)
        do(a, 10)
        do(a, 10)
        do(a, 10)
        rs.Queue(a, 4)
        do(a, 10)
        do(a, 10)
        do(a, 10)
        do(a, 10)
    }()

    sleep(10)

    // B
    go func() {
        rs.Queue(b, 3)
        do(b, 10)
        do(b, 10)
        do(b, 10)
        rs.Queue(b, 1)
        do(b, 10)
    }()

    sleep(10)
    
    // C
    go func() {
        rs.Queue(c, 3)
        do(c, 10)
        do(c, 10)
        do(c, 10)
    }()

    sleep(300)

    rs.Close()

}
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
                rs.Queue(x, 10)
                wg.Done()
                for k := 0; k < 10; k++ {
                    _ = k
                    rs.Proceed(x)
                }
            }(j)
        }
        wg.Wait()
        rs.Close()
    }
}

/*
BenchmarkHoldSwitch_5              15895             90652 ns/op
BenchmarkHoldSwitch_20             10000            320783 ns/op
BenchmarkHoldSwitch_50              4040            867337 ns/op
BenchmarkHoldSwitch_100              608           2215341 ns/op
BenchmarkHoldSwitch_1000               3        1099553892 ns/op

BenchmarkRailSwitch_5              17856             66431 ns/op
BenchmarkRailSwitch_20              4416            270465 ns/op
BenchmarkRailSwitch_50              1714            695491 ns/op
BenchmarkRailSwitch_100              825           1440796 ns/op
BenchmarkRailSwitch_1000              78          18320271 ns/op
*/

func BenchmarkRailSwitch_5(b *testing.B) {
    benchmarkRailSwitch(b.N, 5)
}
func BenchmarkRailSwitch_20(b *testing.B) {
    benchmarkRailSwitch(b.N, 20)
}
func BenchmarkRailSwitch_50(b *testing.B) {
    benchmarkRailSwitch(b.N, 50)
}
func BenchmarkRailSwitch_100(b *testing.B) {
    benchmarkRailSwitch(b.N, 100)
}
func BenchmarkRailSwitch_1000(b *testing.B) {
    benchmarkRailSwitch(b.N, 1000)
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
    do := func(p, s int) {
        t_sleep(s)
        t_time(n[p])
        rs.Proceed(p)
    }

    t_reset()

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
        rs.Queue(a, 1)
        do(a, 10)
    }()

    t_sleep(10)

    // B
    go func() {
        rs.Queue(b, 3)
        do(b, 10)
        do(b, 10)
        do(b, 10)
        rs.Queue(b, 1)
        do(b, 10)
        rs.Queue(b, 2)
        do(b, 10)
        do(b, 10)
    }()

    t_sleep(10)
    
    // C
    go func() {
        rs.Queue(c, 3)
        do(c, 10)
        do(c, 10)
        do(c, 10)
        rs.Queue(c, 5)
        do(c, 10)
        do(c, 10)
        do(c, 10)
        do(c, 10)
        do(c, 10)
    }()

    t_sleep(300)

    rs.Close()

}

func TestRailSwitch2(t *testing.T) {

    rs := NewRailSwitch()

    pr := func(at, i int) {
        fmt.Printf("%c%d,", at + 'b', i)
    }
    do := func(at int) {
        pr(at, 1)
        t_sleep(1)
        pr(at, 2)
        t_sleep(1)
        pr(at, 3)
        t_sleep(1)
        pr(at, 4)
        t_sleep(1)
        rs.Proceed(at)
    }
    repeat := func(at, c int) {
        if rs.Queue(at, c) {
            for i := 0; i < c; i++ {
                go do(at)
            }
        }
    }
    ready := func(at int) {
        char := rune(at) + 'B'
        rs.OnStart(at, func() {
            fmt.Printf("START %c - ", char)
        })
        rs.OnEnd(at, func() {
            fmt.Printf(" END %c\n", char)
        })
    }

    ready(-1)
    ready(0)
    ready(1)
    ready(2)

    go func() {
        repeat(-1, 3)
        t_sleep(2)
        repeat(-1, 3)
        t_sleep(5)
        repeat(-1, 3)
    }()
    go func() {
        repeat(0, 3)
        t_sleep(2)
        repeat(0, 3)
        t_sleep(5)
        repeat(0, 3)
    }()
    t_sleep(1)
    go func() {
        repeat(1, 3)
        t_sleep(7)
        repeat(1, 3)
        t_sleep(2)
        repeat(1, 3)
    }()
    t_sleep(1)
    go func() {
        repeat(2, 3)
        t_sleep(1)
        repeat(2, 3)
        t_sleep(7)
        repeat(2, 3)
    }()

    <- time.After(100 * time.Millisecond)

    rs.Close()

}
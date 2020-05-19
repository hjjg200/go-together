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

/*
BenchmarkHoldSwitch_5              55105             28913 ns/op
BenchmarkHoldSwitch_20             15007             89056 ns/op
BenchmarkHoldSwitch_50             10000            204347 ns/op
BenchmarkHoldSwitch_100             5316            598053 ns/op
BenchmarkHoldSwitch_1000             127          12047837 ns/op

BenchmarkRailSwitch_5              61608             19181 ns/op
BenchmarkRailSwitch_20             16225             74290 ns/op
BenchmarkRailSwitch_50              6300            190849 ns/op
BenchmarkRailSwitch_100             2973            395289 ns/op
BenchmarkRailSwitch_1000             218           5334570 ns/op
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
        rs.Queue(a, 1)
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
        rs.Queue(b, 2)
        do(b, 10)
        do(b, 10)
    }()

    sleep(10)
    
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

    sleep(300)

    rs.Close()

}

func TestRailSwitch2( t *testing.T ) {

    rs := NewRailSwitch()

    rs.OnStart(1, func(){
        log("  A STARTED")
    })
    rs.OnEnd(1, func() {
        log("  A ENDED")
    })

    do := func(name string, at int, ms int) {
        log("ENTERING", name)
        if rs.Queue(at, 1) {
            log("ENTERED", name)
            time.Sleep(time.Millisecond * time.Duration(ms))
            log("EXITING", name)
            rs.Proceed(at)
            log("EXITED", name)
        }
    }

    go do( "A1", 1, 50 )
    time.Sleep(time.Nanosecond)
    go do( "A2", 1, 100 )
    time.Sleep(time.Nanosecond)
    go do( "B1", 2, 200 )
    time.Sleep(time.Nanosecond)
    go do( "C1", 3, 50 )
    time.Sleep(time.Nanosecond)
    go do( "C2", 3, 60 )
    time.Sleep(time.Nanosecond)
    go do( "B2", 2, 50 )
    time.Sleep(time.Nanosecond)

    rs.Close()

}

func TestRailSwitch3(t *testing.T) {

    rs := NewRailSwitch()

    sleep := func() {
        time.Sleep(time.Nanosecond)
    }
    do := func(at int) {
        fmt.Print(1)
        sleep()
        fmt.Print(2)
        sleep()
        fmt.Print(3)
        sleep()
        fmt.Print(4)
        sleep()
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
        char := rune(at) + 'A'
        rs.OnStart(at, func() {
            fmt.Printf("START %c - ", char)
        })
        rs.OnEnd(at, func() {
            fmt.Printf(" END %c\n", char)
        })
    }

    ready(0)
    ready(1)
    ready(2)

    go func() {
        repeat(0, 3)
        repeat(0, 3)
        repeat(0, 3)
    }()
    sleep()
    go func() {
        repeat(1, 3)
        repeat(1, 3)
        repeat(1, 3)
    }()
    sleep()
    go func() {
        repeat(2, 3)
        repeat(2, 3)
        repeat(2, 3)
    }()

    <- time.After(100 * time.Millisecond)

    rs.Close()

}
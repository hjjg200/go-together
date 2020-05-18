package together

import (
    "testing"
    "time"
)


func benchmarkRailSwitch(bn, sz int) {
    for i := 0; i < bn; i++ {
        rs := NewRailSwitch()
        for j := 0; j < sz; j++ {
            go func(x int) {
                rs.Queue(x, 1)
                _ = x
                rs.Proceed(x)
            }(j)
        }
        rs.Wait()
    }
}


func BenchmarkRailSwitch_20(b *testing.B) {
    benchmarkRailSwitch(b.N, 20)
}
func BenchmarkRailSwitch_30(b *testing.B) {
    benchmarkRailSwitch(b.N, 30)
}
func BenchmarkRailSwitch_50(b *testing.B) {
    benchmarkRailSwitch(b.N, 50)
}
func BenchmarkRailSwitch_100(b *testing.B) {
    benchmarkRailSwitch(b.N, 100)
}
func BenchmarkRailSwitch_200(b *testing.B) {
    benchmarkRailSwitch(b.N, 200)
}
func BenchmarkRailSwitch_300(b *testing.B) {
    benchmarkRailSwitch(b.N, 300)
}
func BenchmarkRailSwitch_500(b *testing.B) {
    benchmarkRailSwitch(b.N, 500)
}
func BenchmarkRailSwitch_1000(b *testing.B) {
    benchmarkRailSwitch(b.N, 1000)
}
func BenchmarkRailSwitch_2000(b *testing.B) {
    benchmarkRailSwitch(b.N, 2000)
}
func BenchmarkRailSwitch_5000(b *testing.B) {
    benchmarkRailSwitch(b.N, 5000)
}
func BenchmarkRailSwitch_10000(b *testing.B) {
    benchmarkRailSwitch(b.N, 10000)
}
func BenchmarkRailSwitch_12000(b *testing.B) {
    benchmarkRailSwitch(b.N, 12000)
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
        time.Sleep(time.Nanosecond * time.Duration(s))
    }
    do := func(p, s int) {
        sleep(s)
        t.Logf("%s", n[p])
        rs.Proceed(p)
    }

    // A
    go func() {
        rs.Queue(a, 3)
        do(a, 10)
        do(a, 10)
        do(a, 10)
        rs.Queue(a, 2)
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

    sleep(10)

    <- time.After(time.Nanosecond * 300)

    rs.Wait()

}
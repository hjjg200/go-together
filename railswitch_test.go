package together

import (
    "fmt"
    "testing"
    "time"
)

// BenchmarkHoldSwitch1         200          16711914 ns/op
// BenchmarkRailSwitch1    1000000000               0.211 ns/op

func BenchmarkRailSwitch1(b *testing.B) {
    rs := NewRailSwitch()
    for i := 0; i < 20000; i++ {
        go func(x int) {
            rs.Queue(x, 1)
            _ = x
            rs.Proceed(x)
        }(i)
    }
    rs.Wait()
}

func TestRailSwitch1(t *testing.T) {

    rs := NewRailSwitch()
    sleep := func(s int) {
        time.Sleep(time.Millisecond * time.Duration(s))
    }
    
    const (
        a = iota
        b
        c
    )

    rs.Queue(a, 3)
    go func() {
        sleep(50)
        fmt.Println("a")
        rs.Proceed(a)
        sleep(50)
        fmt.Println("a")
        rs.Proceed(a)
        sleep(50)
        fmt.Println("a")
        rs.Proceed(a)
    }()
    rs.Queue(b, 3)
    go func() {
        sleep(50)
        fmt.Println("b")
        rs.Proceed(b)
        sleep(50)
        fmt.Println("b")
        rs.Proceed(b)
        sleep(50)
        fmt.Println("b")
        rs.Proceed(b)
    }()
    rs.Queue(a, 3)
    go func() {
        sleep(50)
        fmt.Println("a")
        rs.Proceed(a)
        sleep(50)
        fmt.Println("a")
        rs.Proceed(a)
        sleep(50)
        fmt.Println("a")
        rs.Proceed(a)
    }()

    rs.Wait()

    rs.Queue(b, 3)
    go func() {
        sleep(50)
        fmt.Println("b")
        rs.Proceed(b)
        sleep(50)
        fmt.Println("b")
        rs.Proceed(b)
        sleep(50)
        fmt.Println("b")
        rs.Proceed(b)
    }()

    rs.Wait()

}
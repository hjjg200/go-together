package together

import (
    "fmt"
    "testing"
    "time"
    "sync"
)

var start = time.Now()

func log( args ...interface{} ) {
    n := time.Now()
    d := int( n.Sub( start ) / time.Millisecond )
    args = append( []interface{}{ d, "-" }, args... )
    fmt.Println( args... )
}

func benchmarkHoldSwitch(bn, sz int) {
    for i := 0; i < bn; i++ {
        hs := NewHoldSwitch()
        wg := sync.WaitGroup{}
        wg.Add(sz)
        for j := 0; j < sz; j++ {
            go func(x int) {
                hs.Add(x, 1)
                wg.Done()
                _ = x
                hs.Done(x)
            }(j)
        }
        wg.Wait()
        hs.Close()
    }
}

/*
BenchmarkHoldSwitch_5              57049             30108 ns/op
BenchmarkHoldSwitch_20             13810             92816 ns/op
BenchmarkHoldSwitch_50              9967            183696 ns/op
BenchmarkHoldSwitch_100             6855            340940 ns/op
*/

func BenchmarkHoldSwitch_5(b *testing.B) {
    benchmarkHoldSwitch(b.N, 5)
}
func BenchmarkHoldSwitch_20(b *testing.B) {
    benchmarkHoldSwitch(b.N, 20)
}
func BenchmarkHoldSwitch_50(b *testing.B) {
    benchmarkHoldSwitch(b.N, 50)
}
func BenchmarkHoldSwitch_100(b *testing.B) {
    benchmarkHoldSwitch(b.N, 100)
}
func BenchmarkHoldSwitch_1000(b *testing.B) {
    benchmarkHoldSwitch(b.N, 1000)
}

func TestHoldSwitch01( t *testing.T ) {

    hs := NewHoldSwitch()

    hs.Handlers( 1, func() {
        log( "  A STARTED" )
    }, func() {
        log( "  A ENDED" )
    } )

    do := func( name string, at int, ms int ) {
        log( "ENTERING", name )
        hs.Add( at, 1 )
        log( "ENTERED", name )
        time.Sleep( time.Millisecond * time.Duration( ms ) )
        log( "EXITING", name )
        hs.Done( at )
        log( "EXITED", name )
    }

    go do( "A1", 1, 50 )
    time.Sleep( time.Nanosecond )
    go do( "A2", 1, 100 )
    time.Sleep( time.Nanosecond )
    go do( "B1", 2, 200 )
    time.Sleep( time.Nanosecond )
    go do( "C1", 3, 50 )
    time.Sleep( time.Nanosecond )
    go do( "C2", 3, 60 )
    time.Sleep( time.Nanosecond )
    go do( "B2", 2, 50 )
    time.Sleep( time.Nanosecond )

    hs.WaitAll()
    hs.Close()

}
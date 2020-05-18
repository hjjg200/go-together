package together

import (
    "fmt"
    "testing"
    "time"
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
        for j := 0; j < sz; j++ {
            go func(x int) {
                hs.Add(x, 1)
                _ = x
                hs.Done(x)
            }(j)
        }
        hs.Close()
    }
}

/*
BenchmarkHoldSwitch_3              56670             25646 ns/op
BenchmarkHoldSwitch_10             30292             51208 ns/op
BenchmarkHoldSwitch_20             12342             92063 ns/op
BenchmarkHoldSwitch_30             10000            117204 ns/op
BenchmarkHoldSwitch_50              6571            168362 ns/op
BenchmarkHoldSwitch_100             4164            265172 ns/op
*/

func BenchmarkHoldSwitch_3(b *testing.B) {
    benchmarkHoldSwitch(b.N, 3)
}
func BenchmarkHoldSwitch_10(b *testing.B) {
    benchmarkHoldSwitch(b.N, 10)
}
func BenchmarkHoldSwitch_20(b *testing.B) {
    benchmarkHoldSwitch(b.N, 20)
}
func BenchmarkHoldSwitch_30(b *testing.B) {
    benchmarkHoldSwitch(b.N, 30)
}
func BenchmarkHoldSwitch_50(b *testing.B) {
    benchmarkHoldSwitch(b.N, 50)
}
func BenchmarkHoldSwitch_100(b *testing.B) {
    benchmarkHoldSwitch(b.N, 100)
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
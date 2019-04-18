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
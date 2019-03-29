package together

import (
    "testing"
    "sync"
)

func BenchmarkMap01( b *testing.B ) {

    var wg sync.WaitGroup
    wg.Add( 2 )

    mp := NewMap()

    for i := 0; i < 100; i++ {
        mp.Set( i, i )
    }

    go func() {
        for i := 0; i < 100; i++ {
            _, _ = mp.Get( 2 )
        }
        wg.Done()
    }()

    go func() {
        for i := 0; i < 100; i++ {
            mp.Set( 2, 2 )
        }
        wg.Done()
    }()

    wg.Wait()
    mp.hs.Close()

}
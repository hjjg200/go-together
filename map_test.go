package together

import (
    "testing"
    "sync"
)

func BenchmarkMap50K( b *testing.B ) { benchmarkMap( b, 50000 ) }
func BenchmarkMap50K_builtin( b *testing.B ) { benchmarkMap_builtin( b, 50000 ) }
func BenchmarkMap100K( b *testing.B ) { benchmarkMap( b, 100000 ) }
func BenchmarkMap100K_builtin( b *testing.B ) { benchmarkMap_builtin( b, 100000 ) }
func BenchmarkMap200K( b *testing.B ) { benchmarkMap( b, 200000 ) }
func BenchmarkMap200K_builtin( b *testing.B ) { benchmarkMap_builtin( b, 200000 ) }
func BenchmarkMap300K( b *testing.B ) { benchmarkMap( b, 300000 ) }
func BenchmarkMap300K_builtin( b *testing.B ) { benchmarkMap_builtin( b, 300000 ) }
func BenchmarkMap500K( b *testing.B ) { benchmarkMap( b, 500000 ) }
func BenchmarkMap500K_builtin( b *testing.B ) { benchmarkMap_builtin( b, 500000 ) }
func BenchmarkMap1M( b *testing.B ) { benchmarkMap( b, 1000000 ) }
func BenchmarkMap1M_builtin( b *testing.B ) { benchmarkMap_builtin( b, 1000000 ) }

func benchmarkMap( b *testing.B, c int ) {
    var wg sync.WaitGroup
    wg.Add( 2 )

    mp := NewMap()

    for i := 0; i < 100; i++ {
        mp.Set( i, i )
    }

    go func() {
        for i := 0; i < c; i++ {
            _, _ = mp.Get( i % 100 )
        }
        wg.Done()
    }()

    go func() {
        for i := 0; i < c; i++ {
            mp.Set( i % 100, i % 4 )
        }
        wg.Done()
    }()

    wg.Wait()
}

func benchmarkMap_builtin( b *testing.B, c int ) {

    var wg sync.WaitGroup
    wg.Add( 2 )

    mp := make( map[int] int )
    mx := sync.Mutex{}

    for i := 0; i < 100; i++ {
        mp[i] = i
    }

    go func() {
        for i := 0; i < c; i++ {
            mx.Lock()
            _, _ = mp[i % 100]
            mx.Unlock()
        }
        wg.Done()
    }()

    go func() {
        for i := 0; i < c; i++ {
            mx.Lock()
            mp[i % 100] = i % 4
            mx.Unlock()
        }
        wg.Done()
    }()

    wg.Wait()

}
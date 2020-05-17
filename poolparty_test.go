package together

import (
    "runtime"
    "testing"
)

func BenchmarkPoolPartyCPU_0(b *testing.B) {
    pp := NewPoolParty(runtime.NumCPU())
    for i := 0; i < b.N; i++ {
        pp.Join(func() {
            for x := 0; x < 100000; x++ {
                _ = x
            }
        })
    }
    pp.Close()
}

func BenchmarkPoolPartyCPU_1(b *testing.B) {
    pp := NewPoolParty(runtime.NumCPU() + 1)
    for i := 0; i < b.N; i++ {
        pp.Join(func() {
            for x := 0; x < 100000; x++ {
                _ = x
            }
        })
    }
    pp.Close()
}
package together

import (
    "runtime"
    "testing"
)

func benchmarkPoolParty(cpu, bn int) {
    pp := NewPoolParty(cpu)
    for i := 0; i < bn; i++ {
        pp.Join(func() {
            for x := 0; x < 5000000; x++ {
                _ = x
            }
        })
    }
    pp.Close()
}

func BenchmarkPoolPartyCPU_1(b *testing.B) {
    benchmarkPoolParty(1, b.N)
}

func BenchmarkPoolPartyCPU_plus0(b *testing.B) {
    benchmarkPoolParty(runtime.NumCPU(), b.N)
}

func BenchmarkPoolPartyCPU_plus1(b *testing.B) {
    benchmarkPoolParty(runtime.NumCPU() + 1, b.N)
}
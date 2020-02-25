package together

import (
    "testing"
)

func do(pp *PoolParty, i int) {
    pp.Join(func() {
        i++
    })
}

func BenchmarkPoolParty1(t *testing.B) {
    pp := NewPoolParty(4)
    for i := 0; i < 1000000; i++ {
        do(pp, i)
    }
    pp.Close()
}

func BenchmarkPoolParty2(t *testing.B) {
    pp := NewPoolParty(2)
    for i := 0; i < 1000000; i++ {
        do(pp, i)
    }
    pp.Close()
}
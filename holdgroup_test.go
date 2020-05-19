package together

import (
    "testing"
)

// BenchmarkIntMap                  5735191               247 ns/op
// BenchmarkIntInterfaceMap         2162712               728 ns/op

func BenchmarkIntMap(b *testing.B) {
    m := make(map[int] struct{})
    for i := 0; i < b.N; i++ {
        m[i] = struct{}{}
        _ = m[i]
    }
}

func BenchmarkIntInterfaceMap(b *testing.B) {
    m := make(map[interface{}] struct{})
    for i := 0; i < b.N; i++ {
        m[i] = struct{}{}
        _ = m[i]
    }
}

func BenchmarkHoldGroup_int1000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        hg := NewHoldGroup()
        for k := 0; k < 1000; k++ {
            hg.HoldAt(k)
            _ = i
            hg.UnholdAt(k)
        }
    }
}

func BenchmarkHoldGroup_int641000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        hg := NewHoldGroup()
        for k := 0; k < 1000; k++ {
            x := int64(k)
            hg.HoldAt(x)
            _ = i
            hg.UnholdAt(x)
        }
    }
}

func BenchmarkHoldGroup_interface1000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        hg := NewHoldGroup()
        for k := 0; k < 1000; k++ {
            it := float64(k) // HoldGroup does not use float64
            hg.HoldAt(it)
            _ = i
            hg.UnholdAt(it)
        }
    }
}
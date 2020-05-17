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
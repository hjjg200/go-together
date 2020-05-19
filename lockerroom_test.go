package together

import (
    "testing"
)

func TestLockerRoom1(t *testing.T) {
    lr := NewLockerRoom()
    for k := 0; k < 1000; k++ {
        lr.Lock(k)
        _ = k
        lr.Unlock(k)
    }
}

func BenchmarkLockerRoom_int1000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        lr := NewLockerRoom()
        for k := 0; k < 1000; k++ {
            lr.Lock(k)
            _ = i
            lr.Unlock(k)
        }
    }
}

func BenchmarkLockerRoom_int641000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        lr := NewLockerRoom()
        for k := 0; k < 1000; k++ {
            x := int64(k)
            lr.Lock(x)
            _ = i
            lr.Unlock(x)
        }
    }
}

func BenchmarkLockerRoom_interface1000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        lr := NewLockerRoom()
        for k := 0; k < 1000; k++ {
            it := float64(k) // Lockerroom does not use float64
            lr.Lock(it)
            _ = i
            lr.Unlock(it)
        }
    }
}

func BenchmarkLockerRoom_intptr1000(b *testing.B) {
    for i := 0; i < b.N; i++ {
        lr := NewLockerRoom()
        for k := 0; k < 1000; k++ {
            n := k
            lr.Lock(&n)
            _ = i
            lr.Unlock(&n)
        }
    }
}

func BenchmarkIntAsInt(b *testing.B) {
    for i := 0; i < b.N; i++ {
        for k := 0; k < 10000; k++ {
            _ = k
        }
    }
}
func BenchmarkIntAsUint(b *testing.B) {
    for i := 0; i < b.N; i++ {
        for k := 0; k < 10000; k++ {
            _ = uint(k)
        }
    }
}
func BenchmarkIntAsUint64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        for k := 0; k < 10000; k++ {
            _ = uint64(k)
        }
    }
}
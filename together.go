package together

import (
    "runtime"
)

func NumCPU() int {
    return runtime.NumCPU()
}
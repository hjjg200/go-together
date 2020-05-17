package together

import (
    "time"
)

type Passer struct {
    p chan struct{}
    i time.Duration
}

func NewPasser(i time.Duration) *Passer {
    p := make(chan struct{})
    go func() {
        p <- struct{}{}
    }()

    return &Passer{
        p: p,
        i: i,
    }
}

func (ps *Passer) Pass() {
    <- ps.p
    go func() {
        time.Sleep(ps.i)
        ps.p <- struct{}{}
    }()
}

// Deprecated SetInterval as it cannot affect the running timer
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

// Passer defers from time.Ticker in that the first call does not get delayed

// Deprecated SetInterval as it cannot affect the running timer
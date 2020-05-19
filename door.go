package together

import (
    "sync"
    "time"
)

type Door struct {
    c chan struct{}
    i time.Duration
    l time.Time
    m sync.Mutex
}

func NewDoor(i time.Duration) *Door {

    d := &Door{}
    d.c = make(chan struct{})
    d.i = i
    d.l = time.Now()
    
    go func() {
        d.c <- struct{}{}
    }()

    return d

}

func(d *Door) Knock() {

    <- d.c

    go func() {

        d.m.Lock()
        i := d.i
        l := d.l
        d.m.Unlock()

        time.Sleep(i)

        d.m.Lock()
        // Ensure unchanged
        if d.i == i && d.l == l {
            d.c <- struct{}{}
            d.l = time.Now()
        }
        d.m.Unlock()

    }()

}

// Change running timer
// Hang if interval is being changed

func(d *Door) SetInterval(i time.Duration) {

    d.m.Lock()
    if d.i == i {
        d.m.Unlock()
        return
    }

    i0   := d.i
    d.i   = i
    n    := time.Now()
    past := n.Sub(d.l)

    // Pending in channel
    if past >= i0 {
        d.m.Unlock()
        return
    }

    // Past
    if past >= i {
        go func() {
            d.c <- struct{}{}
            d.l = time.Now()
            d.m.Unlock()
        }()
        return
    }

    // Not yet
    left := d.l.Add(i).Sub(n)

    go func() {
        time.Sleep(left)
        d.c <- struct{}{}
        d.l = time.Now()
        d.m.Unlock()
    }()

}
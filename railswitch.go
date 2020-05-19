package together

import (
    "sync"
)

type train struct {
    delta int
    granted chan struct{}
    c1 chan struct{}
}

type rail struct {
    at int
    queue chan *train
}

type RailSwitch struct {
    at int
    value int
    
    rails map[int] *rail
    
    ctrain chan *train
    cat chan int
    queued map[int] bool
    mu sync.Mutex
    rg sync.Mutex
    closed bool
    closer chan struct{}
}

func NewRailSwitch() *RailSwitch {
    rs := new(RailSwitch)
    rs.rails = make(map[int] *rail)
    rs.ctrain = make(chan *train)
    rs.cat = make(chan int)
    rs.at = -1
    rs.queued = make(map[int] bool)

    go func() {
        rs.at = <- rs.cat
        for t := range rs.ctrain {
            rs.value += t.delta
            if rs.value == 0 {
                if rs.at == -1 {
                    t.c1 <- struct{}{}
                    break
                }
                rs.at = <- rs.cat
            }
            t.c1 <- struct{}{}
        }
    }()

    return rs
}

func(rs *RailSwitch) Queue(at, delta int) bool {

    rs.rg.Lock()
    if rs.closed && delta > 0 {
        rs.rg.Unlock()
        return false
    }

    r, ok := rs.rails[at]
    if !ok {
        queue := make(chan *train)
        r = &rail{at, queue}

        rs.rails[at] = r
        rs.queued[at] = false
        go func() {
            for t := range queue {
                if rs.at != r.at && !rs.queued[r.at] {
                    rs.queued[r.at] = true
                    rs.cat <- r.at
                    rs.queued[r.at] = false
                }
                rs.ctrain <- t
                <- t.c1
                t.granted <- struct{}{}
            }
        }()
    }
    rs.rg.Unlock()

    granted := make(chan struct{}, 1)
    c1 := make(chan struct{}, 1)
    r.queue <- &train{
        delta, granted, c1,
    }
    <- granted
    return true

}

func(rs *RailSwitch) Proceed(at int) {
    rs.Queue(at, -1)
}

func(rs *RailSwitch) Close() {
    rs.rg.Lock()
    rs.closed = true
    rs.rg.Unlock()

    rs.Queue(-1, 0)

    close(rs.cat)
    close(rs.ctrain)
    rs.cat = nil
    rs.ctrain = nil
    for _, r := range rs.rails {
        close(r.queue)
    }
    rs.rails = nil
    rs.queued = nil
}
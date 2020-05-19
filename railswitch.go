package together

import (
    "sync"
)

type train struct {
    delta int
    mid chan struct{}
    end chan struct{}
}

type rail struct {
    at int
    queue chan *train
}

type RailSwitch struct {
    at int
    value int
    
    rails    map[int] *rail
    queued   map[int] bool
    cat      chan int
    ctrain   chan *train
    registry sync.Mutex
    closed   bool
}

func NewRailSwitch() *RailSwitch {
    rs := new(RailSwitch)
    rs.at     = -1
    rs.rails  = make(map[int] *rail)
    rs.queued = make(map[int] bool)
    rs.cat    = make(chan int)
    rs.ctrain = make(chan *train)

    go func() {

        rs.at = <- rs.cat

        for t := range rs.ctrain {
            rs.value += t.delta
            if rs.value == 0 {
                if rs.at == -1 {
                    t.mid <- struct{}{}
                    return
                }
                rs.at = <- rs.cat
            }
            t.mid <- struct{}{}
        }

    }()

    return rs
}

func(rs *RailSwitch) Queue(at, delta int) bool {
    if at < 0 {
        panic("together: at must be 0 or higher")
    }
    if delta == 0 {
        panic("together: delta must not be 0")
    }
    return rs.queue(at, delta)
}

func(rs *RailSwitch) queue(at, delta int) bool {

    rs.registry.Lock()
    if rs.closed && delta > 0 {
        rs.registry.Unlock()
        return false
    }

    r, ok := rs.rails[at]
    if !ok {

        queue := make(chan *train)
        r      = &rail{at, queue}

        rs.rails[at]  = r
        rs.queued[at] = false

        go func() {
            for t := range queue {
                if rs.at != r.at && !rs.queued[r.at] {
                    rs.queued[r.at] = true
                    rs.cat <- r.at
                    rs.queued[r.at] = false
                }
                rs.ctrain <- t
                <- t.mid
                t.end <- struct{}{}
            }
        }()
    
    }
    rs.registry.Unlock()

    mid := make(chan struct{}, 1)
    end := make(chan struct{}, 1)
    r.queue <- &train{
        delta, mid, end,
    }
    <- end
    return true

}

func(rs *RailSwitch) Proceed(at int) {
    rs.Queue(at, -1)
}

func(rs *RailSwitch) Close() {

    rs.registry.Lock()
    rs.closed = true
    rs.registry.Unlock()

    rs.queue(-1, 0)

    close(rs.cat)
    close(rs.ctrain)
    rs.cat    = nil
    rs.ctrain = nil
    
    for _, r := range rs.rails {
        close(r.queue)
    }
    rs.rails  = nil
    rs.queued = nil

}
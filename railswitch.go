package together

import (
    "fmt"
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

    sttg map[int] func() // start triggers
    edtg map[int] func() // end triggers
}

func NewRailSwitch() *RailSwitch {
    rs := new(RailSwitch)

    rs.at     = -1
    rs.rails  = make(map[int] *rail)
    rs.queued = make(map[int] bool)
    rs.cat    = make(chan int)
    rs.ctrain = make(chan *train)

    rs.sttg = make(map[int] func())
    rs.edtg = make(map[int] func())

    go func() {

        // First rail
        rs.at = <- rs.cat
        if start := rs.sttg[rs.at]; start != nil { start() }

        for t := range rs.ctrain {

            rs.value += t.delta

            if rs.value == 0 {

                if t.delta == 0 { // closer train
                    t.end <- struct{}{}
                    return
                }

                // For later use
                at0    := rs.at // old at
                end, _ := rs.edtg[at0]

                // Make switch
                rs.at = -1
                t.mid <- struct{}{}
                rs.at = <- rs.cat

                // If rail changed
                if rs.at != at0 {
                    // End of a rail
                    if end != nil { end() }

                    // Start of a rail
                    if start := rs.sttg[rs.at]; start != nil { start() }
                }

                continue

            }

            if rs.value < 0 {
                panic("together: RailSwitch's value must not be below 0")
            }

            t.mid <- struct{}{}
        }

    }()

    return rs
}

func(rs *RailSwitch) Queue(at, delta int) bool {

    if delta <= 0 {
        panic("together: delta must be above 0")
    }

    rs.registry.Lock()
    if rs.closed {
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

    rs.queue(at, delta)
    return true

}

func(rs *RailSwitch) Proceed(at int) {
    if rs.at != at {
        panic("together: proceed attempt for stopped rail")
    }
    rs.queue(at, -1)
}

func(rs *RailSwitch) queue(at, delta int) {
    if at < 0 {
        panic("together: at must not be below 0")
    }

    mid := make(chan struct{}, 1)
    end := make(chan struct{}, 1)
    rs.rails[at].queue <- &train{
        delta, mid, end,
    }
    <- end
}

func(rs *RailSwitch) OnStart(at int, t func()) { rs.sttg[at] = t }
func(rs *RailSwitch) OnEnd(at int, t func())   { rs.edtg[at] = t }

func(rs *RailSwitch) Close() error {

    // Close registry
    rs.registry.Lock()
    if rs.closed {
        rs.registry.Unlock()
        return fmt.Errorf("together: RailSwitch is already closed")
    }

    rs.closed = true

    // send closer train
    rs.cat <- -1 // closer rail
    end := make(chan struct{}, 1)
    rs.ctrain <- &train{
        0, nil, end,
    }
    <- end

    rs.registry.Unlock()

    // Cleanup
    close(rs.cat)
    close(rs.ctrain)
    rs.cat    = nil
    rs.ctrain = nil

    for _, r := range rs.rails {
        close(r.queue)
    }
    rs.rails  = nil
    rs.queued = nil

    return nil

}
package together

import (
    "fmt"
    "sync"
)

type rail struct {
    at int
    queue chan int
}

type RailSwitch struct {
    at int
    value int
    mu sync.Mutex

    rails  chan rail
    queues map[int] chan int
}

/*

Queue 
    // queue has to create new queue when there is no queue for `at`
    // queue has to just pass when current is at `at`
    // queue has to stop when current is not at `at`
    // queue should not use goroutine for sending as it should stop
    // current has to receive in order to let pass

Proceed
    // proceed has to panic when it is not at `at`
    // proceed has to deduct 1 from current queue

*/

func NewRailSwitch() *RailSwitch {

    rs        := new(RailSwitch)
    rs.rails   = make(chan rail)
    rs.queues  = make(map[int] chan int)

    go func() {
        for c := range rs.rails {

            at      := c.at
            queue   := c.queue
            if queue == nil {
                continue
            }

            rs.at    = at
            rs.value = 0

            for c := range queue {
                rs.value += c
                if rs.value == 0 {
                    close(queue)
                    delete(rs.queues, at)
                }
            }

        }
    }()

    return rs

}

func(rs *RailSwitch) Queue(at, delta int) {

    queue, ok := rs.queues[at]
    if !ok {
        rs.mu.Lock()
        _, ok = rs.queues[at]
        if !ok {
            queue = make(chan int)
            rs.queues[at] = queue
            go func() {
                rs.rails <- rail{at, queue}
            }()
        }
        rs.mu.Unlock()
    }

    queue <- delta

}

func(rs *RailSwitch) Proceed(at int) {
    if rs.at != at {
        panic(fmt.Sprintf("together: invalid proceed call for %d while it is at %d", at, rs.at))
    }
    rs.queues[at] <- -1
}

func(rs *RailSwitch) Wait() {
    rs.rails <- rail{}
}

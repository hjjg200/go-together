package together

import (
    "sync"
)

type PoolParty struct {
    pl sync.WaitGroup
    cl sync.WaitGroup

    // The number of threads the PoolParty will use.
    max int

    // The channels of functions the PoolParty will use.
    pool []chan func()
}

func NewPoolParty(threads int) *PoolParty {

    pp := new(PoolParty)

    if threads < 1 {
        threads = 1
    }

    pp.max  = threads
    pp.pool = make([]chan func(), threads)
    for i := range pp.pool {
        pp.pool[i] = make(chan func())
    }

    pp.start()

    return pp

}

func(pp *PoolParty) start() {

    pp.cl.Add(len(pp.pool))
    for i := range pp.pool {
        go func(n int) {
            for f := range pp.pool[n] {
                f()
                pp.pl.Done()
            }
            pp.cl.Done()
        }(i)
    }

}

func(pp *PoolParty) Join(f func()) {
    pp.pl.Add(1)
    go func() {
        mi := pp.minIndex()
        pp.pool[mi] <- f
    }()
}

func(pp *PoolParty) Wait() {
    pp.pl.Wait()
}

func(pp *PoolParty) Close() {

    pp.Wait()

    for i := range pp.pool {
        close(pp.pool[i])
    }

    pp.cl.Wait()
    pp.pool = nil

}

func(pp *PoolParty) minIndex() int {
    mi := 0
    for i := range pp.pool {
        if len(pp.pool[i]) < len(pp.pool[mi]) {
            mi = i
        }
    }
    return mi
}

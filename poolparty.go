package together

import (
    "errors"
    "sync"
)

var (
    ErrThreadCount = errors.New( "together: the given thread count is invalid. The thread count is set to 1." )
)

type PoolParty struct {
    // Mutex is used when assigning functions to channels.
    mx sync.Mutex

    // WaitGroup is used to wait for the current tasks to finish.
    wg sync.WaitGroup

    // The number of threads the PoolParty will use.
    max int

    // The channels of functions the PoolParty will use.
    pool []chan func()
}

func NewPoolParty( threads int ) *PoolParty {

    pp     := new( PoolParty )

    if threads < 1 {
        panic( ErrThreadCount )
        threads = 1
    }

    pp.max  = threads
    pp.pool = make( []chan func(), threads )
    for i := range pp.pool {
        pp.pool[i] = make( chan func() )
    }

    pp.start()

    return pp

}

func ( pp *PoolParty ) start() {

    loop := func( i int ) {
        for f := range pp.pool[i] {
            f()
            pp.wg.Done()
        }
    }

    for i := range pp.pool {
        go loop( i )
    }

}

func ( pp *PoolParty ) Join( f func() ) {

    pp.mx.Lock()

    mi := pp.minIndex()
    pp.wg.Add( 1 )
    pp.pool[mi] <- f

    pp.mx.Unlock()

}

func ( pp *PoolParty ) Wait() {
    pp.wg.Wait()
}

func ( pp *PoolParty ) Close() {

    // Wait for every queue to end
    pp.Wait()

    for i := range pp.pool {
        close( pp.pool[i] )
    }

    pp.pool  = nil

}

func ( pp *PoolParty ) minIndex() int {

    mi := 0
    for i := range pp.pool {
        if len( pp.pool[i] ) < len( pp.pool[mi] ) {
            mi = i
        }
    }
    return mi

}
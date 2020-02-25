package together

import (
    "sync"
)

type HoldGroup struct {
    // Mutex for new key creation
    mx sync.Mutex

    // Map of mutex
    locks map[interface{}] *sync.Mutex

    // Queue of key creation
    queue chan func()
}

func NewHoldGroup() *HoldGroup {

    hg := new(HoldGroup)
    hg.locks = make(map[interface{}] *sync.Mutex)
    hg.queue = make(chan func())

    go func() {
        for f := range hg.queue {
            f()
        }
    }()

    return hg

}

func(hg *HoldGroup) HoldAt(key interface{}) {

    hg.mx.Lock()

    var (
        wg sync.WaitGroup
    )

    wg.Add(1)
    hg.queue <- func() {
        if _, ok := hg.locks[key]; !ok {
            hg.locks[key] = new(sync.Mutex)
        }
        wg.Done()
    }
    wg.Wait()
    hg.mx.Unlock()
    hg.locks[key].Lock()

}

func(hg *HoldGroup) UnholdAt(key interface{}) {
    hg.locks[key].Unlock()
}
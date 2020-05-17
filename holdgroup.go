package together

import (
    "sync"
)

type HoldGroup struct {
    // Mutex for new key creation
    mx sync.Mutex

    // Map of mutex
    locks map[interface{}] *sync.Mutex
}

func NewHoldGroup() *HoldGroup {
    hg := new(HoldGroup)
    hg.locks = make(map[interface{}] *sync.Mutex)

    return hg
}

func(hg *HoldGroup) HoldAt(key interface{}) {

    lock, ok := hg.locks[key]
    if !ok {
        hg.mx.Lock()
        _, ok = hg.locks[key]
        if !ok {
            lock = new(sync.Mutex)
            hg.locks[key] = lock
        }
        hg.mx.Unlock()
    }

    lock.Lock()

}

func(hg *HoldGroup) UnholdAt(key interface{}) {
    hg.locks[key].Unlock()
}
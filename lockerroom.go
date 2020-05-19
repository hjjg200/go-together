package together

import (
    "sync"
    "reflect"
)

type LockerRoom struct {
    mu  sync.Mutex
    im  map[int64] *sync.Mutex
    sm  map[string] *sync.Mutex
    itm map[interface{}] *sync.Mutex
    pm  map[uintptr] *sync.Mutex
}

func NewLockerRoom() *LockerRoom {
    lr := new(LockerRoom)
    lr.im  = make(map[int64] *sync.Mutex)
    lr.sm  = make(map[string] *sync.Mutex)
    lr.itm = make(map[interface{}] *sync.Mutex)
    lr.pm  = make(map[uintptr] *sync.Mutex)
    return lr
}

func(lr *LockerRoom) Lock(key interface{}) {

    var lock *sync.Mutex
    var ok bool

    switch cast := key.(type) {
    case int:    lock, ok = lr.im[int64(cast)]
    case string: lock, ok = lr.sm[cast]
    case int64:  lock, ok = lr.im[cast]
    case uint:   lock, ok = lr.im[int64(cast)]
    case uint64: lock, ok = lr.im[int64(cast)]
    default: 
        v := reflect.ValueOf(key)
        if v.Type().Kind() == reflect.Ptr {
            lock, ok = lr.pm[v.Pointer()]
        } else {
            lock, ok = lr.itm[cast]
        }
    }
    
    if !ok {
        lr.mu.Lock()
        lock = &sync.Mutex{}
        switch cast := key.(type) {
        case int:    lr.im[int64(cast)] = lock
        case string: lr.sm[cast] = lock
        case int64:  lr.im[cast] = lock
        case uint:   lr.im[int64(cast)] = lock
        case uint64: lr.im[int64(cast)] = lock
        default: 
            v := reflect.ValueOf(key)
            if v.Type().Kind() == reflect.Ptr {
                lr.pm[v.Pointer()] = lock
            } else {
                lr.itm[cast] = lock
            }
        }
        lr.mu.Unlock()
    }

    lock.Lock()

}

func(lr *LockerRoom) Unlock(key interface{}) {
    switch cast := key.(type) {
    case int:    lr.im[int64(cast)].Unlock()
    case string: lr.sm[cast].Unlock()
    case int64:  lr.im[cast].Unlock()
    case uint:   lr.im[int64(cast)].Unlock()
    case uint64: lr.im[int64(cast)].Unlock()
    default: 
        v := reflect.ValueOf(key)
        if v.Type().Kind() == reflect.Ptr {
            lr.pm[v.Pointer()].Unlock()
        } else {
            lr.itm[cast].Unlock()
        }
    }
}

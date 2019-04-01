package together

import (
    "sync"
)

type Map struct {
    body map[interface{}] interface{}
    hs *HoldSwitch
    mx sync.Mutex
}

type pair struct {
    Key interface{}
    Value interface{}
}

const (
    c_mapStateRead = iota + 0
    c_mapStateWrite
)

func NewMap() *Map {
    return &Map{
        body: make( map[interface{}] interface{} ),
        hs: NewHoldSwitch(),
    }
}

func( mp *Map ) Len() int {
    mp.hs.Add( c_mapStateRead, 1 )
    length := len( mp.body )
    mp.hs.Done( c_mapStateRead )
    return length
}

func( mp *Map ) Get( key interface{} ) ( interface{}, bool ) {
    mp.hs.Add( c_mapStateRead, 1 )
    val, ok := mp.body[key]
    mp.hs.Done( c_mapStateRead )
    return val, ok
}

func( mp *Map ) Set( key, val interface{} ) {
    mp.hs.Add( c_mapStateWrite, 1 )
    mp.mx.Lock()
    mp.body[key] = val
    mp.mx.Unlock()
    mp.hs.Done( c_mapStateWrite )
}

func( mp *Map ) Values() <- chan interface{} {

    ch := make( chan interface{} )

    go func() {
        mp.hs.Add( c_mapStateRead, 1 )
        for _, v := range mp.body {
            ch <- v
        }
        mp.hs.Done( c_mapStateRead )
        close( ch )
    }()

    return ch

}

func( mp *Map ) Keys() <- chan interface{} {

    ch := make( chan interface{} )

    go func() {
        mp.hs.Add( c_mapStateRead, 1 )
        for k := range mp.body {
            ch <- k
        }
        mp.hs.Done( c_mapStateRead )
        close( ch )
    }()

    return ch

}

func( mp *Map ) Pairs() <- chan pair {

    ch := make( chan pair )

    go func() {
        mp.hs.Add( c_mapStateRead, 1 )
        for k, v := range mp.body {
            ch <- pair{ k, v }
        }
        mp.hs.Done( c_mapStateRead )
        close( ch )
    }()

    return ch

}
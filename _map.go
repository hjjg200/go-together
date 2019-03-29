package together

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
    c_mapStateRead
    c_mapStateWrite
)

func( mp *Map ) Len() int {
    mp.hs.Add( c_mapStateRead, 1 )
    defer mp.hs.Done( c_mapStateRead )
    return len( mp.body )
}

func( mp *Map ) Values() <- chan interface{} {

    ch := make( chan interface{} )

    go func() {
        for _, v := range mp.body {
            ch <- v
        }
        close( ch )
    }()

    return ch

}

func( mp *Map ) Keys() <- chan interface{} {

    ch := make( chan interface{} )

    go func() {
        for k := range mp.body {
            ch <- k
        }
        close( ch )
    }()

    return ch

}

func( mp *Map ) Pairs() <- chan pair {

    ch := make( chan interface{} )

    go func() {
        for k, v := range mp.body {
            ch <- pair{ k, v }
        }
        close( ch )
    }()

    return ch

}
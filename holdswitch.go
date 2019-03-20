package together

import (
    "sync"
)

type HoldSwitch struct {
    at     int
    mx     sync.Mutex
    count  int
    queue  chan wait
    closer chan struct{}

    beginHandlers map[int] func()
    endHandlers   map[int] func()
}

type wait struct {
    at     int
    delta  int
    closer chan struct{}
}

/*

ADD 1 1 - at 1, 1
ADD 2 1 - wait for 1 to end
ADD 1 1 - at 1, 2
ADD 3 1 - wait for 1 to end
ADD 1 -1 - at 1, 1
ADD 1 -1 - at 1, 0 look for the next queue
ADD 2 1 - prolonged - at 2, 1
ADD 3 1 - prolonged - wait for 2 to end
ADD 2 -1, at 2, 0 look for the next queue
ADD 3 1 - prolonged - at 3, 1
ADD 3 -1 - at 3, 0 job done

*/

func NewHoldSwitch() *HoldSwitch {

    hs := &HoldSwitch{
        at: -1,
        count: 0,
        queue: make( chan wait ),
        closer: nil,

        beginHandlers: make( map[int] func() ),
        endHandlers: make( map[int] func() ),
    }

    go hs.loop()
    return hs

}

func( hs *HoldSwitch ) loop() {

    aside := make( []wait, 0 )
    do    := func( w wait ) {
        if hs.at == w.at {
            hs.count += w.delta
            w.closer <- struct{}{}
        } else {
            if hs.count == 0 {

                // At
                lastAt  := hs.at
                hs.at    = w.at
                hs.count = w.delta

                // Begin and End
                if lastAt != -1 {
                    if ed, ok := hs.endHandlers[lastAt]; ok && ed != nil {
                        ed()
                    }
                }
                if bg, ok := hs.beginHandlers[w.at]; ok && bg != nil {
                    bg()
                }

                // Close
                w.closer <- struct{}{}

            } else {
                aside = append( aside, w )
            }
        }
    }

    for {
        select {
        case w := <- hs.queue:

            hs.mx.Lock()

            // Main
            do( w )

            // Aside
            if len( aside ) > 0 {
                cp := make( []wait, len( aside ) )
                copy( cp, aside )
                aside = make( []wait, 0 )
                for _, aw := range cp {
                    do( aw )
                }
            }

            // Call the closer if any
                // Not wait and count is at 0
            if len( aside ) == 0 && len( hs.queue ) == 0 && hs.count == 0 {
                if hs.closer != nil {
                    hs.closer <- struct{}{}
                    hs.closer = nil
                }
            }

            hs.mx.Unlock()

        }
    }

}

func( hs *HoldSwitch ) Add( at, delta int ) {

    if at < 0 {
        // panic
        return
    }

    //
    closer := make( chan struct{}, 1 )

    hs.mx.Lock()
    hs.queue <- wait{
        at: at,
        delta: delta,
        closer: closer,
    }
    hs.mx.Unlock()

    <- closer

}

func( hs *HoldSwitch ) Done( at int ) {
    hs.Add( at, -1 )
}

func( hs *HoldSwitch ) WaitAll() {
    hs.closer = make( chan struct{}, 1 )
    <- hs.closer
}

func( hs *HoldSwitch ) Handlers( at int, begin func(), end func() ) {
    hs.beginHandlers[at] = begin
    hs.endHandlers[at]   = end
}
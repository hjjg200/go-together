package together


type train struct {
    at int
    delta int
    granted chan struct{}
}

type RailSwitch struct {
    at int
    value int
    
    buffered []*train
    queue chan *train
    open bool
}

func NewRailSwitch() *RailSwitch {
    rs := new(RailSwitch)
    rs.buffered = make([]*train, 0)
    rs.queue    = make(chan *train)

    go func() {
        Loop:
        for {

            buf := make([]*train, len(rs.buffered))
            copy(buf, rs.buffered)
            rs.buffered = make([]*train, 0)

            for _, t := range buf {
                rs.register(t)
                print(rs.open); print(t.at); print(" -b "); println(t.delta)

                if !rs.open { continue Loop }
            }

            for t := range rs.queue {
                rs.register(t)

                print(rs.open); print(t.at); print(" -  "); println(t.delta)
                if !rs.open { continue Loop }
            }

        }
    }()

    return rs
}

func bufferHelper(b []*train) {
    for _, t :=range b {
        print("B"); print(t.at); print("+"); println(t.delta)
    }
}

func(rs *RailSwitch) register(t *train) {

    if !rs.open {
        rs.at = t.at
        rs.open = true
    }

    if t.at == rs.at {
        rs.value += t.delta
        if rs.value == 0 {
            rs.open = false
        }
        if t.granted != nil {
            t.granted <- struct{}{}
        }
    } else {
        rs.buffered = append(rs.buffered, t)
    }
}

func(rs *RailSwitch) Queue(at, delta int) {
    granted := make(chan struct{}, 1)
    rs.queue <- &train{
        at, delta, granted,
    }
    <- granted
}

func(rs *RailSwitch) Proceed(at int) {
    rs.Queue(at, -1)
}

func(rs *RailSwitch) Wait() {
    // add panic for at < 0
    rs.Queue(-1, 0)
}
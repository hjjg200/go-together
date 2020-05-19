package together

import (
    "fmt"
    "testing"
    "time"
)

func TestDoor1(t *testing.T) {

    st := time.Now()
    log := func() {
        fmt.Printf("%v\n", time.Now().Sub(st))
    }

    d := NewDoor(time.Millisecond * 60)

    d.Knock()
    log()

    d.Knock()
    log()

    d.SetInterval(time.Millisecond * 20)

    d.Knock()
    log()

    d.Knock()
    log()

    go func() {
        d.SetInterval(time.Millisecond * 40)
        d.SetInterval(time.Millisecond * 20)
        d.SetInterval(time.Millisecond * 60)
    }()

    d.Knock()
    log()
    
    d.Knock()
    log()

    d.Knock()
    log()
    
    d.Knock()
    log()

}

func TestDoor2(t *testing.T) {
    st := time.Now()
    log := func() {
        fmt.Printf("%v\n", time.Now().Sub(st))
    }

    d := NewDoor(time.Millisecond * 60)

    go func() {
        d2 := NewDoor(time.Millisecond * 20)
        d2.Knock()
        d.SetInterval(time.Millisecond * 30)
        d2.Knock()
        d.SetInterval(time.Millisecond * 20)
        d2.Knock()
        d.SetInterval(time.Millisecond * 50)
        d2.Knock()
        d.SetInterval(time.Millisecond * 20)
    }()

    d.Knock()
    log()
    d.Knock()
    log()
    d.Knock()
    log()
    d.Knock()
    log()
    d.Knock()
    log()

}
package together

import (
    "testing"
    "time"
)

func TestDoor1(t *testing.T) {

    t_reset()

    d := NewDoor(time.Millisecond * 60)

    d.Knock()
    t_time()

    d.Knock()
    t_time()

    d.Set(time.Millisecond * 20)

    d.Knock()
    t_time()

    d.Knock()
    t_time()

    go func() {
        d.Set(time.Millisecond * 40)
        d.Set(time.Millisecond * 20)
        d.Set(time.Millisecond * 60)
    }()

    d.Knock()
    t_time()
    
    d.Knock()
    t_time()

    d.Knock()
    t_time()
    
    d.Knock()
    t_time()

}

func TestDoor2(t *testing.T) {
    
    t_reset()

    d := NewDoor(time.Millisecond * 60)

    go func() {
        d2 := NewDoor(time.Millisecond * 20)
        d2.Knock()
        d.Set(time.Millisecond * 30)
        d2.Knock()
        d.Set(time.Millisecond * 20)
        d2.Knock()
        d.Set(time.Millisecond * 50)
        d2.Knock()
        d.Set(time.Millisecond * 20)
    }()

    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()

}

func TestDoor3(t *testing.T) {

    t_reset()
    
    d := NewDoor(time.Millisecond * 60)

    go func() {
        d.Set(time.Millisecond * 50)
        d.Set(time.Millisecond * 30)
        d.Set(time.Millisecond * 20)
        d.Set(time.Millisecond * 10)
    }()
    
    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()
    d.Knock()
    t_time()

}

func TestDoor4(t *testing.T) {
    
    t_reset()
    d := NewDoor(time.Second * 1)
    for i := 2; i < 6; i++ {
        d.Knock()
        t_time("!")
        d.Set(time.Second * time.Duration(i))
    }

}
# package together

```go
import "github.com/hjjg200/go/together"
```

Package together provides concurrent utilities. It is currently in its beta stage and there might be bugs.

## Index

* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [type LockerRoom](#LockerRoom)
    * [func NewLockerRoom() *LockerRoom](#NewLockerRoom)
    * [func(lr *LockerRoom) Lock(key interface{})](#LockerRoom.Lock)
    * [func(lr *LockerRoom) Unlock(key interface{})](#LockerRoom.Unlock)
* [type RailSwitch](#RailSwitch)
    * [func NewRailSwitch() *RailSwitch](#NewRailSwitch)
    * [func(rs *RailSwitch) Queue(at, delta int) bool](#RailSwitch.Queue)
    * [func(rs *RailSwitch) Proceed(at int)](#RailSwitch.Proceed)
    * [func(rs *RailSwitch) OnStart(at int, t func())](#RailSwitch.OnStart)
    * [func(rs *RailSwitch) OnEnd(at int, t func())](#RailSwitch.OnEnd)
    * [func(rs *RailSwitch) Close() error](#RailSwitch.Close)
* [type Door](#Door)
    * [func NewDoor(i time.Duration) *Door](#NewDoor)
    * [func(d *Door) Knock()](#Door.Knock)
    * [func(d *Door) Set(i time.Duration)](#Door.Set)


## <a name="pkg-variables" href="#">Constants</a>

```go
const (
)
```

## <a name="pkg-variables" href="#">Variables</a>

```go
var (
)
```

## <a name="LockerRoom" href="#">type LockerRoom</a>

```go
type LockerRoom struct {
    // contains filtered or unexported fields
}
```

LockerRoom is a collective mutex group.

### <a name="NewLockerRoom" href="#">func NewLockerRoom</a>

```go
func NewLockerRoom() *LockerRoom
```

Returns a new LockerRoom.

### <a name="LockerRoom.HoldAt" href="#">func(*LockerRoom) Lock</a>

```go
func(lr *LockerRoom) Lock(key interface{})
```

Locks the mutex for the given key. If it is the first time to lock the key, the LockerRoom assigns a mutex to the relevant map according to the type of the key. `int`, `uint`, `int64`, and `uint64` are assigned to `map[int64] *sync.Mutex`; `string` is assigned to `map[string]`; pointers such as `*int` and `*struct{}` are assigned to `map[uintptr]`; the others are all assigned to `map[interface{}]`. A LockerRoom has separate maps for increased performance.

### <a name="LockerRoom.Unlock" href="#">func(*LockerRoom) Unlock</a>

```go
func(lr *LockerRoom) Unlock(key interface{})
```

Attempts to unlock the mutex for the given key.


## <a name="RailSwitch" href="#">type RailSwitch</a>

```go
type RailSwitch struct {
    // contains filtered or unexported fields
}
```

Like a rail switch that allows only one rail to proceed, RailSwitch allows only one group to operate while others are waiting for their turn. Turn is determined on first come, first served basis; the goroutine that called `RailSwitch.Queue` first will receive the next turn.

### <a name="NewRailSwitch" href="#">func NewRailSwitch</a>

```go
func NewRailSwitch() *RailSwitch
```

Creates a new RailSwitch. You need to assign triggers before calling `RailSwitch.Queue` in order to ensure they are triggered at the first call for queue.

### <a name="RailSwitch.Queue" href="#">func(*RailSwitch) Queue</a>

```go
func(rs *RailSwitch) Queue(at, delta int) bool
```

Queue hangs until its group gains turn and returns true if the group successfully gained turn and false if the RailSwitch was closed. Therefore, you need to wrap it with a if block for proper use.

And number below 0 is not allowed for `at`. They are reserved for internal use.

```go
if rs.Queue(groupCleanup, 1) {
    // Successfully gained turn
}
// Failed to gain turn
```

`for` blocks can be used for repeating tasks:

```go
// Repeating a task every minute
for rs.Queue(groupMain, 1) {
    foo()
    rs.Proceed(groupMain)
    time.Sleep(1 * time.Minute)
}
```


### <a name="RailSwitch.Proceed" href="#">func(*RailSwitch) Proceed</a>

```go
func(rs *RailSwitch) Proceed(at int)
```

Proceed notifies the RailSwitch that one of the current group has completed its task. And when the value -- figuratively, the number of remaining trains -- of the current "rail" reaches 0, the RailSwitch gives turn to the awaiting group.

It panics if it is called for a group that is not in turn.

### <a name="RailSwitch.OnStart" href="#">func(*RailSwitch) OnStart</a>

```go
func(rs *RailSwitch) OnStart(at int, t func())
```

OnStart is used to set trigger for a certain group. The trigger is guaranteed to start and complete prior to the first operation of that group.

### <a name="RailSwitch.OnEnd" href="#">func(*RailSwitch) OnEnd</a>

```go
func(rs *RailSwitch) OnEnd(at int, t func())
```

OnEnd is used to set trigger for a certain group. The trigger is guaranteed to start and complete when a switch happens or the RailSwitch is closed; if it is when a switch happened the trigger is run and completed prior to the start trigger and the first operation of the next group.

### <a name="RailSwitch.Close" href="#">func(*RailSwitch) Close</a>

```go
func(rs *RailSwitch) Close() error
```

Close gracefully closes the RailSwitch. It waits for the "trains" and "rails" that were already queued for operating, and once they are done, it closes underlying channels and internal goroutines of RailSwitch and blocks the trains and rails came after the call for Close and blocks any consecutive train.

It returns an error if it is already closed.


## <a name="Door" href="#">type Door</a>

```go
type Door struct {
    // contains filtered or unexported fields
}
```

Door is similar to `time.Ticker` but defers from it in that it does not delay the first call and that you can change the interval for running timer unlike `time.Ticker.Reset`.

### <a name="NewDoor" href="#">func NewDoor</a>

```go
func NewDoor(i time.Duration) *Door
```

Creates a door for the given duration.

### <a name="Door.Knock" href="#">func(*Door) Knock</a>

```go
func(d *Door) Knock()
```

Knock does not hang if it is the first call and if it has passed interval amount of time since its last call. It hangs if it has not passed interval amount of time since its last call.

### <a name="Door.Set" href="#">func(*Door) Set</a>

```go
func(d *Door) Set(i time.Duration)
```

Set sets a new duration for the Door. It replaces the timer that Knock may be waiting for. And it hangs if there is any duration change ongoing. Therefore, it is highly likely to create a deadlock if it is called twice in the same goroutine.

```go
d := NewDoor(time.Second * 1)
for i := 2; i < 7; i++ {
    d.Knock()
    fmt.Println("!")
    d.Set(time.Second * time.Duration(i))
}

// T+00.000000 - !
// T+02.000000 - !
// T+05.000000 - !
// T+09.000000 - !
// T+14.000000 - !
```
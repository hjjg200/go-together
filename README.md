# package together

** Unstable prototype **

```go
import "github.com/hjjg200/go/together"
```

Package together provides helpful features for concurrent programming.

## Index

* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func NumCPU() int](#NumCPU)
* [type HoldGroup](#HoldGroup)
    * [func NewHoldGroup() *HoldGroup](#NewHoldGroup)
    * [func ( hg *HoldGroup ) HoldAt( key interface{} )](#HoldGroup.HoldAt)
    * [func ( hg *HoldGroup ) UnholdAt( key interface{} )](#HoldGroup.UnholdAt)
* [type HoldSwitch](#HoldSwitch)
    * [func NewHoldSwitch() *HoldSwitch](#NewHoldSwitch)
    * [func ( hs *HoldSwitch ) AddA( delta int )](#HoldSwitch.AddA)
    * [func ( hs *HoldSwitch ) AddB( delta int )](#HoldSwitch.AddB)
    * [func ( hs *HoldSwitch ) DoneA()](#HoldSwitch.DoneA)
    * [func ( hs *HoldSwitch ) DoneB()](#HoldSwitch.DoneB)
* [type PoolParty](#PoolParty)
    * [func NewPoolParty( threads int ) ( *PoolParty, error )](#NewPoolParty)
    * [func ( pp *PoolParty ) Join( f func() )](#PoolParty.Join)
    * [func ( pp *PoolParty ) Wait()](#PoolParty.Wait)
    * [func ( pp *PoolParty ) Close()](#PoolParty.Close)

## <a name="pkg-variables" href="#pkg-variables">Variables</a>

```go
var (
    ErrThreadCount = errors.New( "together: the given thread count is invalid. The thread count is set to 1." )
)
```

## <a name="NumCPU" href="#NumCPU">func NumCPU</a>

```go
func NumCPU() int
```

Returns `runtime.NumCPU`.

## <a name="HoldGroup" href="#HoldGroup">type HoldGroup</a>

```go
type HoldGroup struct {
    // no exported members
}
```

HoldGroup is a type that can be used as a collective mutex group.

### <a name="NewHoldGroup" href="#NewHoldGroup">func NewHoldGroup</a>

```go
func NewHoldGroup() *HoldGroup
```

Returns a new HoldGroup.

### <a name="HoldGroup.HoldAt" href="#HoldGroup.HoldAt">func ( *HoldGroup ) HoldAt</a>

```go
func ( hg *HoldGroup ) HoldAt( key interface{} )
```

Locks the mutex of the given key.

### <a name="HoldGroup.UnholdAt" href="#HoldGroup.UnholdAt">func ( *HoldGroup ) UnholdAt</a>

```go
func ( hg *HoldGroup ) UnholdAt( key interface{} )
```

Unlocks the mutex of the given key.

## <a name="PoolParty" href="#PoolParty">type PoolParty</a>

```go
type PoolParty struct {
    // no exported members
}
```

PoolParty processes the given queues in number of goroutines. It assigns new queues to the pool that has the least queues. Future plan is to make it assign new queues to the pool that finishes the job.

### <a name="NewPoolParty" href="#NewPoolParty">func NewPoolParty</a>

```go
func NewPoolParty( threads int ) *PoolParty
```

Creates a new PoolParty. It starts a loop that loops through the channels. Use NumCPU to use as many threads as there are. Panics if the given thread count is below 1.

### <a name="PoolParty.Join" href="#PoolParty.Join">func ( *PoolParty ) Join</a>

```go
func ( pp *PoolParty ) Join( f func() )
```

Put f into the channel that has least queues.

### <a nane="PoolParty.Wait" href="#PoolParty.Wait">func ( *PoolParty ) Wait</a>

```go
func ( pp *PoolParty ) Wait()
```

Wait until all the tasks to be finishied.

### <a name="PoolParty.Close" href="#PoolParty.Close">func ( *PoolParty ) Close</a>

```go
func ( pp *PoolParty ) Close()
```

Waits for the PoolParty to finish all the queues and closes the channels.
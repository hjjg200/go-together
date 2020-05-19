# package together

** Unstable prototype **


```import "github.com/hjjg200/go/together"```

Package together provides helpful features for concurrent programming.

## Index

* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* func NumCPU() int
* type HoldGroup
    * func NewHoldGroup() *HoldGroup
    * func ( hg *HoldGroup ) HoldAt( key interface{} )
    * func ( hg *HoldGroup ) UnholdAt( key interface{} )
* type HoldSwitch
    * func NewHoldSwitch() *HoldSwitch
    * func ( hs *HoldSwitch ) AddA( delta int )
    * func ( hs *HoldSwitch ) AddB( delta int )
    * func ( hs *HoldSwitch ) DoneA()
    * func ( hs *HoldSwitch ) DoneB()
    * func ( hs *HoldSwitch ) ReadyA()
    * func ( hs *HoldSwitch ) ReadyB()
* type PoolParty
    * func NewPoolParty( threads int ) ( *PoolParty, error )
    * func ( pp *PoolParty ) start()
    * func ( pp *PoolParty ) Join( f func() )
    * func ( pp *PoolParty ) Wait()
    * func ( pp *PoolParty ) Close()
    * func ( pp *PoolParty ) minIndex() int

## Variables

```go
var (
    ErrThreadCount = errors.New( "together: the given thread count is invalid" )
)
```

## type PoolParty

```go
type PoolParty struct {
    // Mutex is used when assigning functions to channels.
    mx sync.Mutex

    // WaitGroup is used to wait for the current tasks to finish.
    wg sync.WaitGroup

    // The number of threads the PoolParty will use.
    max int

    // The channels of functions the PoolParty will use.
    pool []chan func()

    // The count of queued functions.
    count []int
}
```

### func NewPoolParty

```go
func NewPoolParty( threads int ) ( *PoolParty, error )
```

Creates a new PoolParty. It starts a loop that loops through the channels. Use NumCPU to use as many threads as there are.

### func ( *PoolParty ) start

```go
func ( pp *PoolParty ) start()
```

Starts loops that loop through channels.

### func ( *PoolParty ) Join

```go
func ( pp *PoolParty ) Join( f func() )
```

Put f into the channel that has least queues.

### func ( *PoolParty ) Wait

```go
func ( pp *PoolParty ) Wait()
```

Wait until all the tasks to be finishied.

### func ( *PoolParty ) Close

```go
func ( pp *PoolParty ) Close()
```

Closes the PoolParty that will be no longer used.

### func ( *PoolParty ) minIndex

```go
func ( pp *PoolParty ) minIndex() int
```

Returns the index of the channel that has the least queues.
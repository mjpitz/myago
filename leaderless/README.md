# leaderless

Package leaderless implements leader election without the need for coordination.
It does this using a highly stable set of peers and a hash ring. Since we know
that each node shares the same view of the world, then we know that the computed
hash ring will be the same between instances. However, if peer knowledge is not
stable (for example, members come and go freely) then leaderless can result in a
split brain state where some nodes share a different view of the world until the
other nodes "catch up".

This package is loosely inspired by Uber's ringpop system which seems like it's
used quite extensively.

```go
import github.com/mjpitz/myago/leaderless
```

## Usage

#### type Director

```go
type Director struct {
}
```

Director contains logic for routing requests to a leader or set of replicas.

#### func New

```go
func New() *Director
```

New returns a Director that can aid in the coordination of work within a
cluster.

#### func (\*Director) GetLeader

```go
func (d *Director) GetLeader(key string) (string, bool)
```

GetLeader returns the current "leader" for a given key.

#### func (\*Director) GetReplicas

```go
func (d *Director) GetReplicas(key string, replicas int) ([]string, bool)
```

GetReplicas returns a list of peers to replicate information to given a key.

#### func (\*Director) Start

```go
func (d *Director) Start(ctx context.Context, membership *cluster.Membership) error
```

Start begins the director by observing membership changes in the cluster.

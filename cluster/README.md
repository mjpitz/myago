# cluster
--
    import "github.com/mjpitz/myago/cluster"

Package cluster provides code to manage cluster Membership. Membership can
currently be managed statically (via explicit configuration of active) or
dynamically (using HashiCorp's Serf project).

## Usage

#### type CancelWatch

```go
type CancelWatch func()
```

CancelWatch is used to remove a watch from the cluster membership.

#### type Cluster

```go
type Cluster struct {
}
```

Cluster handles the discovery and management of cluster members. It uses
HashiCorp's Serf and MemberList projects to discover and track active given a
join address.

#### func  New

```go
func New(opts ...Option) *Cluster
```
New constructs a cluster given the provided options.

#### func (*Cluster) Membership

```go
func (c *Cluster) Membership() *Membership
```
Membership returns the underlying membership of the cluster. Useful for
obtaining a snapshot or for manipulating the entries used by the cluster during
testing.

#### func (*Cluster) Start

```go
func (c *Cluster) Start(ctx context.Context) error
```
Start initializes and starts up the cluster. It uses errgroup to spin up the
discovery thread and blocks until the parent context is cancelled or one of the
grouped functions returns an error.

#### type Discovery

```go
type Discovery interface {
	// Start begins the process of filling the membership with active members. This method should block until the
	// context is cancelled.
	Start(ctx context.Context, membership *Membership) error
}
```

Discovery is a generic interface used to manage the membership pool of the
cluster.

#### type GossipDiscovery

```go
type GossipDiscovery struct {
	JoinAddress string
	Config      *serf.Config
}
```

GossipDiscovery uses HashiCorp's Serf library to discover nodes within the
cluster. It requires both TCP and UDP communication to be available.

#### func (*GossipDiscovery) Start

```go
func (g *GossipDiscovery) Start(ctx context.Context, membership *Membership) error
```

#### type Membership

```go
type Membership struct {
}
```

Membership tacks a current list of active within the cluster. It can be
populated manually (useful for testing) or using common discovery mechanisms.

#### func (*Membership) Add

```go
func (m *Membership) Add(peers []string)
```
Add inserts the provided active into the cluster's active list. Operation should
be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) + len(m.left)`.

#### func (*Membership) Left

```go
func (m *Membership) Left(peers []string)
```
Left allows peers to temporarily leave the cluster, but still be considered part
of active membership. Operation should be `O( m log(n) )` where `m = len(peers)`
and `n = len(m.active) + len(m.left)`.

#### func (*Membership) Remove

```go
func (m *Membership) Remove(peers []string)
```
Remove deletes the provided active from the cluster's peer list. Operation
should be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) +
len(m.left)`.

#### func (*Membership) Snapshot

```go
func (m *Membership) Snapshot() ([]string, int)
```
Snapshot returns a copy of the current peer list.

#### func (*Membership) Watch

```go
func (m *Membership) Watch() (<-chan MembershipChange, CancelWatch)
```
Watch allows others to observe changes in the cluster membership.

#### type MembershipChange

```go
type MembershipChange struct {
	Active  []string
	Left    []string
	Removed []string
}
```

MembershipChange describes how the cluster membership has changed to outside
observers.

#### type NoDiscovery

```go
type NoDiscovery struct {
	Peers []string
}
```

NoDiscovery uses a statically provided list of peers to fill Membership.

#### func (*NoDiscovery) Start

```go
func (n *NoDiscovery) Start(ctx context.Context, membership *Membership) error
```

#### type Option

```go
type Option func(cluster *Cluster)
```

Option defines how callers can customize aspects of the cluster.

#### func  WithDiscovery

```go
func WithDiscovery(discovery Discovery) Option
```
WithDiscovery allows alternative peer discovery mechanisms to be plugged in.

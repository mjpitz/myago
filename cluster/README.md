# cluster

Package cluster provides code to manage cluster Membership. Membership can
currently be managed statically (via explicit configuration of active) or
dynamically (using HashiCorp's Serf project).

```go
import github.com/mjpitz/myago/cluster
```

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

#### func New

```go
func New(opts ...Option) *Cluster
```

New constructs a cluster given the provided options.

#### func (\*Cluster) Membership

```go
func (c *Cluster) Membership() *Membership
```

Membership returns the underlying membership of the cluster. Useful for
obtaining a snapshot or for manipulating the entries used by the cluster during
testing.

#### func (\*Cluster) Start

```go
func (c *Cluster) Start(ctx context.Context) error
```

Start initializes and starts up the cluster.

#### type Config

```go
type Config struct {
	NoDiscovery
	GossipDiscovery
}
```

Config provides a common configuration structure for forming clusters. Either
through a list of known addresses (peers) or using gossip to form pool
dynamically.

#### func (\*Config) Start

```go
func (c *Config) Start(ctx context.Context, membership *Membership) error
```

Start controls which discovery mechanism is invoked based on the provided
configuration.

#### type Discovery

```go
type Discovery interface {
	// Start runs the discovery process. Implementations should block, regardless if they're filling or subscribing to
	// the membership pool.
	Start(ctx context.Context, membership *Membership) error
}
```

Discovery provides an abstraction that allows implementers to fill or discover
changes to the underlying membership pool. For example, GossipDiscovery fills
the membership pool with members found via HashiCorp's Serf implementation. The
leaderless.Director package implements this interface to learn about changes in
the underlying membership pool.

#### type GossipDiscovery

```go
type GossipDiscovery struct {
	JoinAddress string       `json:"join_address" usage:"create a cluster dynamically through a single join address"`
	Config      *serf.Config `json:"-"`
}
```

GossipDiscovery uses HashiCorp's Serf library to discover nodes within the
cluster. It requires both TCP and UDP communication to be available.

#### func (\*GossipDiscovery) Start

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

#### func (\*Membership) Add

```go
func (m *Membership) Add(peers []string)
```

Add inserts the provided active into the cluster's active list. Operation should
be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) + len(m.left)`.

#### func (\*Membership) Left

```go
func (m *Membership) Left(peers []string)
```

Left allows peers to temporarily leave the cluster, but still be considered part
of active membership. Operation should be `O( m log(n) )` where `m = len(peers)`
and `n = len(m.active) + len(m.left)`.

#### func (\*Membership) Majority

```go
func (m *Membership) Majority() int
```

Majority computes a cluster majority. This returns a simple majority for the
cluster.

#### func (\*Membership) Remove

```go
func (m *Membership) Remove(peers []string)
```

Remove deletes the provided active from the cluster's peer list. Operation
should be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) + len(m.left)`.

#### func (\*Membership) Snapshot

```go
func (m *Membership) Snapshot() ([]string, int)
```

Snapshot returns a copy of the current peer list.

#### func (\*Membership) Watch

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
	Peers []string `json:"peers" usage:"create a cluster using a static list of addresses"`
}
```

NoDiscovery uses a statically provided list of peers to fill Membership.

#### func (\*NoDiscovery) Start

```go
func (n *NoDiscovery) Start(ctx context.Context, membership *Membership) error
```

#### type Option

```go
type Option func(cluster *Cluster)
```

Option defines how callers can customize aspects of the cluster.

#### func WithDiscovery

```go
func WithDiscovery(discovery Discovery) Option
```

WithDiscovery allows alternative peer discovery mechanisms to be plugged in.

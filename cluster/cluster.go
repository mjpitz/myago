package cluster

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Option defines how callers can customize aspects of the cluster.
type Option func(cluster *Cluster)

// WithDiscovery allows alternative peer discovery mechanisms to be plugged in.
func WithDiscovery(discovery Discovery) Option {
	return func(cluster *Cluster) {
		cluster.discovery = append(cluster.discovery, discovery)
	}
}

// New constructs a cluster given the provided options.
func New(opts ...Option) *Cluster {
	cluster := &Cluster{
		membership: new(Membership),
	}

	for _, opt := range opts {
		opt(cluster)
	}

	return cluster
}

// Cluster handles the discovery and management of cluster members. It uses HashiCorp's Serf and MemberList projects to
// discover and track active given a join address.
type Cluster struct {
	discovery  []Discovery
	membership *Membership
}

// Membership returns the underlying membership of the cluster. Useful for obtaining a snapshot or for manipulating the
// entries used by the cluster during testing.
func (c *Cluster) Membership() *Membership {
	return c.membership
}

// Start initializes and starts up the cluster.
func (c *Cluster) Start(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	submitToGroup := func(discovery Discovery) {
		group.Go(func() error {
			return discovery.Start(ctx, c.membership)
		})
	}

	for _, discovery := range c.discovery {
		submitToGroup(discovery)
	}

	return group.Wait()
}

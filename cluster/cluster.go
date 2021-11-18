// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

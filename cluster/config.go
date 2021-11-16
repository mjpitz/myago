package cluster

import (
	"context"
)

// Config provides a common configuration structure for forming clusters. Either through a list of known addresses
// (peers) or using gossip to form pool dynamically.
type Config struct {
	NoDiscovery
	GossipDiscovery
}

// Start controls which discovery mechanism is invoked based on the provided configuration.
func (c *Config) Start(ctx context.Context, membership *Membership) error {
	switch {
	case len(c.NoDiscovery.Peers) > 0:
		return c.NoDiscovery.Start(ctx, membership)
	case len(c.GossipDiscovery.JoinAddress) > 0:
		return c.GossipDiscovery.Start(ctx, membership)
	}

	return nil
}

package cluster

import (
	"context"
)

// NoDiscovery uses a statically provided list of peers to fill Membership.
type NoDiscovery struct {
	Peers []string
}

func (n *NoDiscovery) Start(ctx context.Context, membership *Membership) error {
	membership.Add(n.Peers)
	<-ctx.Done()
	return nil
}

var _ Discovery = &GossipDiscovery{}

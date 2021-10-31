package cluster

import "context"

// Discovery provides an abstraction that allows implementers to fill or discover changes to the underlying membership
// pool. For example, GossipDiscovery fills the membership pool with members found via HashiCorp's Serf implementation.
// The leaderless.Director package implements this interface to learn about changes in the underlying membership pool.
type Discovery interface {
	// Start runs the discovery process. Implementations should block, regardless if they're filling or subscribing to
	// the membership pool.
	Start(ctx context.Context, membership *Membership) error
}

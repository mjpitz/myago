package cluster

import "context"

// Discovery is a generic interface used to manage the membership pool of the cluster.
type Discovery interface {
	// Start begins the process of filling the membership with active members. This method should block until the
	// context is cancelled.
	Start(ctx context.Context, membership *Membership) error
}

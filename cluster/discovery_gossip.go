package cluster

import (
	"context"

	"github.com/hashicorp/serf/serf"
	"golang.org/x/sync/errgroup"
)

// GossipDiscovery uses HashiCorp's Serf library to discover nodes within the cluster. It requires both TCP and UDP
// communication to be available.
type GossipDiscovery struct {
	JoinAddress string
	Config      *serf.Config
}

func (g *GossipDiscovery) Start(ctx context.Context, membership *Membership) error {
	eventCh := make(chan serf.Event, 16)
	defer close(eventCh)

	g.Config.EventCh = eventCh
	// g.Config.Logger = serfLog
	// g.Config.MemberlistConfig.Logger = serfLog

	serfClient, err := serf.Create(g.Config)
	if err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		defer func() {
			_ = serfClient.Shutdown()
		}()

		_, err := serfClient.Join([]string{g.JoinAddress}, false)
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			case ev := <-eventCh:
				if memberEvent, ok := ev.(serf.MemberEvent); ok {
					updatedPeers := make([]string, 0, len(memberEvent.Members))
					for _, member := range memberEvent.Members {
						updatedPeers = append(updatedPeers, member.Addr.String())
					}

					switch memberEvent.EventType() {
					case serf.EventMemberJoin:
						membership.Add(updatedPeers)
					case serf.EventMemberLeave, serf.EventMemberFailed:
						membership.Left(updatedPeers)
					case serf.EventMemberReap:
						membership.Remove(updatedPeers)
					}
				}
			}
		}
	})

	return group.Wait()
}

var _ Discovery = &GossipDiscovery{}

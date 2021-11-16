package cluster

import (
	"context"

	"github.com/hashicorp/serf/serf"

	"github.com/mjpitz/myago/zaputil"
)

// GossipDiscovery uses HashiCorp's Serf library to discover nodes within the cluster. It requires both TCP and UDP
// communication to be available.
type GossipDiscovery struct {
	JoinAddress string       `json:"join_address" usage:"create a cluster dynamically through a single join address"`
	Config      *serf.Config `json:"-"`
}

func (g *GossipDiscovery) consume(ctx context.Context, membership *Membership, eventCh chan serf.Event) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
}

func (g *GossipDiscovery) Start(ctx context.Context, membership *Membership) error {
	eventCh := make(chan serf.Event, 16)
	defer close(eventCh)

	logger := zaputil.HashicorpStdLogger(zaputil.Extract(ctx))

	g.Config.EventCh = eventCh
	g.Config.Logger = logger
	g.Config.MemberlistConfig.Logger = logger

	serfClient, err := serf.Create(g.Config)
	if err != nil {
		return err
	}

	defer func() {
		_ = serfClient.Shutdown()
	}()

	_, err = serfClient.Join([]string{g.JoinAddress}, false)
	if err != nil {
		return err
	}

	return g.consume(ctx, membership, eventCh)
}

var _ Discovery = &GossipDiscovery{}

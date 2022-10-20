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

	"github.com/hashicorp/serf/serf"
	"go.pitz.tech/lib/logger"
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

	log := logger.HashiCorpStdLogger(logger.Extract(ctx))

	g.Config.EventCh = eventCh
	g.Config.Logger = log
	g.Config.MemberlistConfig.Logger = log

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

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
)

// NoDiscovery uses a statically provided list of peers to fill Membership.
type NoDiscovery struct {
	Peers []string `json:"peers" usage:"create a cluster using a static list of addresses"`
}

func (n *NoDiscovery) Start(ctx context.Context, membership *Membership) error {
	membership.Add(n.Peers)
	<-ctx.Done()

	return nil
}

var _ Discovery = &GossipDiscovery{}

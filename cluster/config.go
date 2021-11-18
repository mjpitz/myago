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

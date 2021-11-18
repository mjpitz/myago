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

import "context"

// Discovery provides an abstraction that allows implementers to fill or discover changes to the underlying membership
// pool. For example, GossipDiscovery fills the membership pool with members found via HashiCorp's Serf implementation.
// The leaderless.Director package implements this interface to learn about changes in the underlying membership pool.
type Discovery interface {
	// Start runs the discovery process. Implementations should block, regardless if they're filling or subscribing to
	// the membership pool.
	Start(ctx context.Context, membership *Membership) error
}

// Copyright (C) 2022 Mya Pitzeruse
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
	"net"
	"time"

	"go.pitz.tech/lib/clocks"
)

// DNSDiscovery uses DNS to resolve cluster membership. Currently, this implementation uses the default DNS resolver
// that comes with Go. I know that the serf library uses something beyond the default implementation, so it might be
// worth exploring this later on.
type DNSDiscovery struct {
	Name            string        `json:"dns_name" usage:"specify the dns name to resolve"`
	ResolveInterval time.Duration `json:"dns_resolve_interval" usage:"how frequently the dns name should be resolved" default:"30s"`
}

func (dns *DNSDiscovery) Start(ctx context.Context, membership *Membership) error {
	left := make(map[string]bool)
	last := make(map[string]bool)

	clock := clocks.Extract(ctx)
	ticker := clock.NewTicker(1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.Chan():
			ticker.Stop()
			ticker = clock.NewTicker(dns.ResolveInterval)

			addrs, err := net.DefaultResolver.LookupIPAddr(ctx, dns.Name)
			if err != nil {
				continue
			}

			add := make([]string, 0, len(addrs))
			leave := make([]string, 0, len(addrs))
			remove := make([]string, 0, len(addrs))

			next := make(map[string]bool, len(addrs))
			for _, addr := range addrs {
				peer := addr.String()
				next[peer] = true

				if !last[peer] {
					add = append(add, peer)
				}
			}

			for peer := range left {
				if !next[peer] {
					remove = append(remove, peer)
				}
			}

			left = make(map[string]bool)

			for peer := range last {
				if !next[peer] {
					left[peer] = true
					leave = append(leave, peer)
				}
			}

			last = next

			if len(add) > 0 {
				membership.Add(add)
			}

			if len(leave) > 0 {
				membership.Left(leave)
			}

			if len(remove) > 0 {
				membership.Remove(remove)
			}
		}
	}
}

var _ Discovery = &DNSDiscovery{}

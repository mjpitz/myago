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

package cluster_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"go.pitz.tech/lib/cluster"
)

func testHarness(t *testing.T, discovery cluster.Discovery, length int) {
	t.Helper()

	c := cluster.New(
		cluster.WithDiscovery(discovery),
	)

	ctx, cancel := context.WithCancel(context.Background())

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return c.Start(ctx)
	})

	group.Go(func() error {
		defer cancel()

		ch, unwatch := c.Membership().Watch()
		defer unwatch()

		var change cluster.MembershipChange
		select {
		case <-ctx.Done():
			return ctx.Err()
		case change = <-ch:
		}

		if len(change.Active) == 0 {
			// skip initial snapshot of state, but it can be a race
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change = <-ch:
			}
		}

		require.Len(t, change.Active, length)
		require.Len(t, change.Left, 0)
		require.Len(t, change.Removed, 0)

		return nil
	})

	err := group.Wait()

	switch {
	case errors.Is(err, context.Canceled):
	case err != nil:
		require.Fail(t, err.Error())
	}
}

func TestConfigGossipDiscovery(t *testing.T) {
	t.Parallel()

	testHarness(t, &cluster.Config{
		GossipDiscovery: cluster.GossipDiscovery{
			JoinAddress: "localhost:7946",
			Config:      serf.DefaultConfig(),
		},
	}, 1)
}

func TestConfigNoDiscovery(t *testing.T) {
	t.Parallel()

	testHarness(t, &cluster.Config{
		NoDiscovery: cluster.NoDiscovery{
			Peers: []string{
				"peer-1",
				"peer-2",
				"peer-3",
			},
		},
	}, 3)
}

func TestConfigDNSDiscovery(t *testing.T) {
	t.Parallel()

	testHarness(t, &cluster.Config{
		DNSDiscovery: cluster.DNSDiscovery{
			Name:            "go.pitz.tech",
			ResolveInterval: 30 * time.Second,
		},
	}, 2)
}

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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/cluster"
)

func TestMembership(t *testing.T) {
	t.Parallel()

	membership := &cluster.Membership{}

	{
		peers, _ := membership.Snapshot()
		require.Len(t, peers, 0)
	}

	allHosts := []string{"host-1", "host-2", "host-3"}
	shuffled := []string{"host-3", "host-1", "host-2"}

	membership.Add(shuffled)

	{
		peers, _ := membership.Snapshot()
		require.Len(t, peers, 3)
		require.Equal(t, allHosts, peers)
	}

	membership.Left([]string{"host-2"})
	{
		peers, n := membership.Snapshot()
		require.Len(t, peers, 3)
		require.Equal(t, 2, n)
		require.Equal(t, []string{"host-1", "host-3", "host-2"}, peers)
	}

	membership.Remove([]string{"host-2"})
	{
		peers, _ := membership.Snapshot()
		require.Len(t, peers, 2)
		require.Equal(t, []string{"host-1", "host-3"}, peers)
	}

	membership.Left(allHosts)
	{
		// all peers should be left peers
		peers, n := membership.Snapshot()
		require.Len(t, peers[:n], 0)
		require.Len(t, peers, 2)
	}

	membership.Remove(allHosts)
	{
		// all peers should be left peers
		peers, _ := membership.Snapshot()
		require.Len(t, peers, 0)
	}
}

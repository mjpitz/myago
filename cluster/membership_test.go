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

package leaderless_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/cluster"
	"github.com/mjpitz/myago/leaderless"
)

func TestDirector(t *testing.T) {
	ctx := context.Background()

	membership := new(cluster.Membership)
	membership.Add([]string{"host-3", "host-1", "host-2"})

	director := leaderless.New()

	{
		timeoutCtx, _ := context.WithTimeout(ctx, time.Second)
		err := director.Start(timeoutCtx, membership)
		require.NoError(t, err)
	}

	leader, ok := director.GetLeader("leader")

	require.True(t, ok)
	require.Equal(t, "host-3", leader)
}

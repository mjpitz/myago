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

package leaderless_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/cluster"
	"go.pitz.tech/lib/leaderless"
)

func TestDirector(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	membership := new(cluster.Membership)
	membership.Add([]string{"host-3", "host-1", "host-2"})

	director := leaderless.New()

	{
		timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		err := director.Start(timeoutCtx, membership)
		require.NoError(t, err)
	}

	leader, ok := director.GetLeader("leader")

	require.True(t, ok)
	require.Equal(t, "host-3", leader)
}

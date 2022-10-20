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

package leaderless

import (
	"context"
	"sync"

	"github.com/dgryski/go-farm"
	"github.com/serialx/hashring"

	"go.pitz.tech/lib/cluster"
)

// New returns a Director that can aid in the coordination of work within a cluster.
func New() *Director {
	return &Director{
		hashRing: hashring.NewWithHash([]string{}, func(data []byte) hashring.HashKey {
			low, high := farm.Hash128(data)

			return &hashring.Int64PairHashKey{
				High: int64(high),
				Low:  int64(low),
			}
		}),
	}
}

// Director contains logic for routing requests to a leader or set of replicas.
type Director struct {
	mu       sync.RWMutex
	hashRing *hashring.HashRing
}

// GetLeader returns the current "leader" for a given key.
func (d *Director) GetLeader(key string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.hashRing.GetNode(key)
}

// GetReplicas returns a list of peers to replicate information to given a key.
func (d *Director) GetReplicas(key string, replicas int) ([]string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// when replicas > pool size, return the full pool
	if size := d.hashRing.Size(); replicas > size {
		replicas = size
	}

	return d.hashRing.GetNodes(key, replicas)
}

func (d *Director) handleMembershipChange(change cluster.MembershipChange) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch {
	case len(change.Active) > 0:
		for _, added := range change.Active {
			d.hashRing = d.hashRing.AddNode(added)
		}
	case len(change.Removed) > 0:
		for _, removed := range change.Removed {
			d.hashRing = d.hashRing.RemoveNode(removed)
		}
	}
}

// Start begins the director by observing membership changes in the cluster.
func (d *Director) Start(ctx context.Context, membership *cluster.Membership) error {
	watch, cancel := membership.Watch()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil
		case change := <-watch:
			d.handleMembershipChange(change)
		}
	}
}

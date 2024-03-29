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

package paxos

import (
	"time"

	"github.com/jonboulle/clockwork"
)

// IDGenerator defines an interface for generating IDs used internally by paxos.
type IDGenerator interface {
	Next() (uint64, error)
}

// ServerIDGenerator returns an IDGenerator that creates an ID using a provided serverID and clock. It works by
// taking a millisecond level timestamp, shifting its value left 8 bits, and or it with the server ID. The leading
// byte can be used to expand this representation later on.
//
//	const (
//		wordView  = 0x0000000000000000
//
//		nowMillis = 0x0000017c96370c09
//		shifted   = 0x00017c96370c0900
//		withSID   = 0x00017c96370c09XX
//	)
//
// As you can see, there is plenty of space for the IDGenerator to function. Obviously, there are limitations with this
// implementation.
//
//  1. 256 max possible instances
//  2. Throughput constrained to 1 op/ms
//
// Granted, some of these aren't _huge_ issues for the types of systems that this _could_ help build.
func ServerIDGenerator(serverID uint8, clock clockwork.Clock) IDGenerator {
	return &serverIDGenerator{
		serverID: serverID,
		clock:    clock,
	}
}

type serverIDGenerator struct {
	serverID uint8
	clock    clockwork.Clock
}

func (s *serverIDGenerator) Next() (uint64, error) {
	millisecond := s.clock.Now().UnixNano() / int64(time.Millisecond)
	id := (uint64(millisecond) << 8) | uint64(s.serverID)

	return id, nil
}

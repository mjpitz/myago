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

// Vote is an internal structure used by multiple components to cast votes on behalf of the acceptor that they're
// communicating with.
type Vote struct {
	// Member contains which member of the cluster cast this vote.
	Member string
	// Payload contains the payload of the message we're voting on. This is usually a Promise or Proposal.
	Payload interface{}
}

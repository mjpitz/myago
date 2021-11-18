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

/*
Package paxos implements the paxos algorithm. The logic is mostly ported from mjpitz/paxos, but with a few
modifications. First, I didn't continue using gRPC as the transport as I wanted something a bit less cumbersome. I've
tried to break down the interface in such a way where different transports _could_ be plugged in. More on that later.

This package is (and likely will be for a while) a work in progress. As it stands, it _should_ support simple paxos.
*/
package paxos

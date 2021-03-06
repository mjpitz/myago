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
Package leaderless implements leader election without the need for coordination. It does this using a highly stable set
of peers and a hash ring. Since we know that each node shares the same view of the world, then we know that the computed
hash ring will be the same between instances. However, if peer knowledge is not stable (for example, members come and go
freely) then leaderless can result in a split brain state where some nodes share a different view of the world until the
other nodes "catch up".

This package is loosely inspired by Uber's ringpop system which seems like it's used quite extensively.
*/
package leaderless

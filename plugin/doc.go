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

// Package plugin provides a simple plugin interface by forking processes and using their stdout/stdin to enable
// communication between the parent process (main-component) and the child (plugin). This is inspired by how protoc and
// its various plugins work. Applications can read arguments, flags, and environment variables provided to the program
// to configure its behaviour, but then stream data from stdin to issue RPCs and write their responses to stdout.
package plugin

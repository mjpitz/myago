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
Package dirset provides discovery of common application directories for things like caching, locks, and logs.
*/
package dirset

// DirectorySet defines a common set of paths that applications can use for a variety of reasons.
type DirectorySet struct {
	CacheDir      string
	StateDir      string
	LocalStateDir string
	LockDir       string
	LogDir        string
}

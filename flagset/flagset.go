// Copyright (C) 2022 Mya Pitzeruse
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

package flagset

import (
	"github.com/urfave/cli/v2"
)

// FlagSet provides additional functionality on top of a collection of flags.
type FlagSet []cli.Flag

// Filter returns a new FlagSet that contains flags allowed by the provided Filter.
func (flags FlagSet) Filter(allow Filter) FlagSet {
	next := FlagSet{}

	for _, flag := range flags {
		if allow(flag) {
			next = append(next, flag)
		}
	}

	return next
}

// Filter allows the user to inspect the flag to determine if it should be in the resulting FlagSet.
type Filter func(flag cli.Flag) bool

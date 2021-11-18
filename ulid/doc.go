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
Package ulid provides code for generating variable length unique, lexigraphic identifiers (ULID) with programmable
fills. Currently, there is a RandomGenerator that can be used to generate ULIDs with a randomized payload. To provide
a custom payload, simply extend the BaseGenerator, and override the Generate method. It's important to call the
BaseGenerator's Generate method, otherwise the skew and timestamp bits won't be set properly.

Unlike the canonical [ULID](https://github.com/ulid/spec), this version holds a placeholder byte for major clock skews
which can often occur in distributed systems. The wire format is as follows: `[ skew ][ millis ][ payload ]`

 - `skew` - 1 byte used to handle major clock skews (reserved, unused)
 - `millis` - 6 bytes of a unix timestamp (should give us until the year 10k or so)
 - `payload` - N bytes for the payload
*/
package ulid

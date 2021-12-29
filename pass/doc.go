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

// Package pass provides password derivation functions backing solutions like Spectre. There are three steps in the
// process. First, you need to derive an Identity key based your name and primary password. This key is unique to
// you (assuming name / password combinations are unique). The second step is to generate a SiteKey. This key is unique
// to you for the site that you're generating the key for. Finally, the last step is to generate a SitePassword using
// the derived SiteKey and associated password format.
package pass

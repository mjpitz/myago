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

package lifecycle

import (
	"context"
)

var systemLifeCycle = &LifeCycle{}

// Defer will enqueue a function that will be invoked by Resolve.
func Defer(fn func(ctx context.Context)) {
	systemLifeCycle.Defer(fn)
}

// Resolve will process all functions that have been enqueued by Defer up until this point.
func Resolve(ctx context.Context) {
	systemLifeCycle.Resolve(ctx)
}

// Setup initializes a shutdown hook that cancels the underlying context.
func Setup(ctx context.Context) context.Context {
	return systemLifeCycle.Setup(ctx)
}

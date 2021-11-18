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

package clocks

import (
	"context"

	"github.com/jonboulle/clockwork"

	"github.com/mjpitz/myago"
)

var contextKey = myago.ContextKey("clocks")

var defaultClock = clockwork.NewRealClock()

// Extract pulls the clock from the provided context. If no clock is found, then the defaultClock is returned.
func Extract(ctx context.Context) clockwork.Clock {
	clock := ctx.Value(contextKey)
	if clock == nil {
		return defaultClock
	}

	return clock.(clockwork.Clock)
}

// ToContext sets the clock on the provided context.
func ToContext(ctx context.Context, clock clockwork.Clock) context.Context {
	return context.WithValue(ctx, contextKey, clock)
}

// Setup sets the defaultClock on the provided context. This can always be overridden later.
func Setup(ctx context.Context) context.Context {
	return ToContext(ctx, defaultClock)
}

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
	"os/signal"
	"sync"
	"syscall"
)

// LifeCycle hooks into various lifecycle events. It allows functions to be deferred en masse.
type LifeCycle struct {
	once     sync.Once
	mu       sync.Mutex
	funcs    []func(ctx context.Context)
	shutdown context.CancelFunc
	halt     context.CancelFunc
}

// Defer will enqueue a function that will be invoked by Resolve.
func (lc *LifeCycle) Defer(fn func(ctx context.Context)) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.funcs = append(lc.funcs, fn)
}

// Resolve will process all functions that have been enqueued by Defer up until this point.
func (lc *LifeCycle) Resolve(ctx context.Context) {
	fns := func() []func(ctx context.Context) {
		lc.mu.Lock()
		defer lc.mu.Unlock()

		fn := append([]func(ctx context.Context){}, lc.funcs...)
		lc.funcs = lc.funcs[len(fn):]

		return fn
	}()

	for i := len(fns); i > 0; i-- {
		fns[i-1](ctx)
	}
}

// Setup initializes a shutdown hook that cancels the underlying context.
func (lc *LifeCycle) Setup(ctx context.Context) context.Context {
	lc.once.Do(func() {
		lc.mu.Lock()
		defer lc.mu.Unlock()

		ctx, lc.shutdown = context.WithCancel(ctx)
		ctx, lc.halt = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	})

	return ctx
}

func (lc *LifeCycle) Shutdown(ctx context.Context) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if lc.shutdown != nil {
		lc.shutdown()
		lc.halt()

		lc.shutdown = nil
		lc.halt = nil
		lc.once = sync.Once{}
	}
}

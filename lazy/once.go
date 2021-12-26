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

package lazy

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// Once will attempt to load a value until one is loaded.
type Once struct {
	// Loader is a function that returns an object and optional error. It conditionally accepts a context value.
	Loader interface{}
	mu     sync.Mutex
	value  interface{}
}

// Get returns the loaded value if set or an error should one occur.
func (o *Once) Get(ctx context.Context) (interface{}, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.value != nil {
		return o.value, nil
	}

	if o.Loader == nil {
		return nil, fmt.Errorf("loader not set")
	}

	v := reflect.ValueOf(o.Loader)
	t := v.Type()
	if v.Kind() != reflect.Func {
		return nil, fmt.Errorf("loader is not a function")
	} else if out := t.NumOut(); out > 2 || out == 0 {
		return nil, fmt.Errorf("incorrect number of return values: %d. expected 1 or 2", out)
	}

	req := []reflect.Value{
		reflect.ValueOf(ctx),
	}

	req = req[:t.NumIn()]

	resp := v.Call(req)
	resp = resp[:t.NumOut()]

	if len(resp) > 1 && !resp[1].IsNil() {
		err, ok := resp[1].Interface().(error)
		if !ok {
			return nil, fmt.Errorf("second return value was not an error")
		} else if err != nil {
			return nil, err
		}
	}

	o.value = resp[0].Interface()
	return o.value, nil
}

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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// Slice provides a configuration file (YAML, JSON, TOML) friendly variant of the cli.StringSlice value. It also
// supports parsing JSON configuration from environment variables.
type Slice[T any] []T

// Set appends the provided value to the underlying slice. If the value appears to be a JSON string, then we attempt to
// unmarshal the JSON before appending the resulting values to the underlying slice. Currently, this mechanism only
// supports primitive structures.
func (s *Slice[T]) Set(value string) (err error) {
	if value == "" {
		return nil
	}

	newValues := make(Slice[T], 0)

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		err = json.Unmarshal([]byte(value), &newValues)
	} else {
		dst := new(T)
		err = parse(value, dst)
		if err == nil {
			newValues = append(newValues, *dst)
		}
	}

	if err != nil {
		return err
	}

	*s = append(*s, newValues...)

	return nil
}

// String serializes the Slice as a json string.
func (s *Slice[T]) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}

	value, _ := json.Marshal(*s)
	return string(value)
}

var _ cli.Generic = &Slice[any]{}

// parse attempts to convert the provided string to the target destination.
func parse(value string, dst any) (err error) {
	switch v := dst.(type) {
	case *string:
		*v = value
	case *bool:
		*v, err = strconv.ParseBool(value)
		if err != nil {
			return err
		}
	case *int, *int8, *int16, *int32, *int64:
		var x int64
		x, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		switch v := v.(type) {
		case *int:
			*v = int(x)
		case *int8:
			*v = int8(x)
		case *int16:
			*v = int16(x)
		case *int32:
			*v = int32(x)
		case *int64:
			*v = x
		}
	case *uint, *uint8, *uint16, *uint32, *uint64:
		var x uint64
		x, err = strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}

		switch v := v.(type) {
		case *uint:
			*v = uint(x)
		case *uint8:
			*v = uint8(x)
		case *uint16:
			*v = uint16(x)
		case *uint32:
			*v = uint32(x)
		case *uint64:
			*v = x
		}
	case *float32, *float64:
		var x float64
		x, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		switch v := v.(type) {
		case *float32:
			*v = float32(x)
		case *float64:
			*v = x
		}
	case *complex64, *complex128:
		var x complex128
		x, err = strconv.ParseComplex(value, 128)
		if err != nil {
			return err
		}

		switch v := v.(type) {
		case *complex64:
			*v = complex64(x)
		case *complex128:
			*v = x
		}
	default:
		return fmt.Errorf("provided destination is not a pointer to a primative")
	}

	return nil
}

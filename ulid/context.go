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

package ulid

import (
	"context"
	"os"
	"strconv"

	"go.pitz.tech/lib/libctx"
)

var (
	contextKey      = libctx.Key("ulid.generator")
	systemGenerator *Generator
)

//nolint:gochecknoinits
func init() {
	skew := byte(1)
	if skewEnv := os.Getenv("MYAGO_ULID_SKEW"); skewEnv != "" {
		s, err := strconv.Atoi(skewEnv)
		if err != nil {
			// not a fan
			panic(err)
		}
		skew = byte(s)
	}

	systemGenerator = NewGenerator(skew, RandomFill)
}

// Extract is used to obtain the generator from a context. If none is present, the system generator is used.
func Extract(ctx context.Context) *Generator {
	val := ctx.Value(contextKey)
	if val == nil {
		return systemGenerator
	}

	return val.(*Generator)
}

// ToContext appends the provided generator to the provided context.
func ToContext(ctx context.Context, generator *Generator) context.Context {
	return context.WithValue(ctx, contextKey, generator)
}

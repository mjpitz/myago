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

package zaputil

import (
	"context"

	"go.uber.org/zap"

	"github.com/mjpitz/myago"
)

var contextKey = myago.ContextKey("zap")

var defaultLogger = zap.NewNop()

// Extract pulls the logger from the provided context. If no logger is found, then the defaultLogger is returned.
func Extract(ctx context.Context) *zap.Logger {
	log := ctx.Value(contextKey)
	if log == nil {
		return defaultLogger
	}

	return log.(*zap.Logger)
}

// ToContext sets the logger on the provided context.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

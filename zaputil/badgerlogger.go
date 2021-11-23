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
	"fmt"

	"go.uber.org/zap"
)

// Logger is an interface pulled from the badger library. It defines the functionality needed by the badger system to
// log messages. It supports a variety of levels and works similar to the fmt.Printf method.
type Logger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

// BadgerLogger wraps the provided logger so badger can log using zap.
func BadgerLogger(log *zap.Logger) Logger {
	return &logger{
		log: log.WithOptions(zap.AddCallerSkip(2)),
	}
}

type logger struct {
	log *zap.Logger
}

func (l *logger) Errorf(s string, i ...interface{}) {
	l.log.Error(fmt.Sprintf(s, i...))
}

func (l *logger) Warningf(s string, i ...interface{}) {
	l.log.Warn(fmt.Sprintf(s, i...))
}

func (l *logger) Infof(s string, i ...interface{}) {
	l.log.Info(fmt.Sprintf(s, i...))
}

func (l *logger) Debugf(s string, i ...interface{}) {
	l.log.Debug(fmt.Sprintf(s, i...))
}

var _ Logger = &logger{}

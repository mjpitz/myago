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

package index

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseDSN string `json:"database_dsn" usage:"specify the connection string for database" default:"file:db.sqlite"`
}

func Open(cfg Config) (*Index, error) {
	var dialector gorm.Dialector

	switch {
	case cfg.DatabaseDSN == "":
		return nil, fmt.Errorf("missing database_dsn flag")
	case strings.HasPrefix(cfg.DatabaseDSN, "postgres:"):
		dialector = postgres.Open(cfg.DatabaseDSN)
	case strings.HasPrefix(cfg.DatabaseDSN, "file:"):
		dialector = sqlite.Open(cfg.DatabaseDSN)
	default:
		return nil, fmt.Errorf("unrecognized database target")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Index{
		db: db,
	}, nil
}

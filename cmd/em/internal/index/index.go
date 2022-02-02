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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Index struct {
	db *gorm.DB
}

func (i *Index) Migrate(schema ...interface{}) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		for _, s := range schema {
			if err := tx.AutoMigrate(s); err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *Index) Index(docs ...interface{}) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Clauses(clause.OnConflict{DoNothing: true})

		for _, doc := range docs {
			tx.Create(doc)
		}

		return nil
	})
}

func (i *Index) Close() error {
	db, err := i.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

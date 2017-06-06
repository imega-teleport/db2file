package mysql

import "database/sql"

type storage struct {
	db *sql.DB
}

func NewStorage(sqlDB *sql.DB) *storage {
	return &storage{
		db: sqlDB,
	}
}

func (s *storage) Groups() ([]Group, err error) {
	return
}
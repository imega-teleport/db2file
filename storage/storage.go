package storage

import (
	"database/sql"
)

type Storage interface {
	GetGroups(out chan<- interface{}, e chan<- error)
	GetProducts(out chan<- interface{}, e chan<- error)
	GetProductsGroups(out chan<- interface{}, e chan<- error)
	GetProductsProperties(out chan<- interface{}, e chan<- error, condition []string)
	GetProductsPropertiesSpecial(out chan<- interface{}, e chan<- error, condition []string)
	CheckCompleteAllTasks() (bool, error)
}

type storage struct {
	db    *sql.DB
	limit int
}

func NewStorage(db *sql.DB, limit int) Storage {
	return &storage{
		db:    db,
		limit: limit,
	}
}

func (s *storage) CheckCompleteAllTasks() (bool, error) {
	var completeTask = 0
	row := s.db.QueryRow("select count(*) as compl from tasks where val = 1")
	if err := row.Scan(&completeTask); err != nil {
		return false, err
	}

	if completeTask == 2 {
		return true, nil
	}

	return false, nil
}

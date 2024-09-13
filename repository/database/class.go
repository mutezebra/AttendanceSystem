package database

import "database/sql"

type ClassRepository struct {
	db *sql.DB
}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{
		db: _db,
	}
}

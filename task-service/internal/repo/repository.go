package repo

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// TODO: ошибки???

type Repository struct {
	db *sql.DB
}

// TODO: изменить параметры
func NewDb(config string) (*Repository, error) {
	db, err := sql.Open("postgres", config)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

package multiarmed_bandit

import (
	"database/sql"
	"fmt"
)

type sqlStorage struct {
	db *sql.DB
}

func NewSqlStorage(db *sql.DB) {
	return &sqlStorage{
		db: db,
	}
}

func (s *sqlStorage) Find(name string) err {
	rows, err := s.db.Query(
		`SELECT e.name, e.gamma, v.weight, v.shows, v.rewards
		FROM experiment e INNER JOIN variant v ON v.experiment_id = v.id
		WHERE name = %s`,
		name,
	)

	if err != nil {
		return err
	}

	variants := []Variant{}
	for rows.Next() {

	}

}

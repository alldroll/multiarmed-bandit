package multiarmed_bandit

import (
	"database/sql"
	"fmt"
)

type sqlStorage struct {
	db *sql.DB
}

func NewSqlStorage(db *sql.DB) *sqlStorage {
	return &sqlStorage{
		db: db,
	}
}

//
func (s *sqlStorage) Find(name string) (Experiment, error) {
	rows, err := s.db.Query(
		`SELECT e.gamma, v.shows, v.rewards
		FROM experiment e INNER JOIN variant v ON v.experiment_name = e.name
		WHERE e.name = %s`,
		name,
	)

	if err != nil {
		return nil, err
	}

	var (
		gamma          float32
		shows, rewards uint32
		variants       = []Variant{}
	)

	for rows.Next() {
		rows.Scan(&gamma, &shows, &rewards)
		variants = append(variants, NewVariant(0, rewards, shows))
	}

	if len(variants) == 0 {
		return nil, nil
	}

	return NewExperiment(name, gamma, variants), nil
}

//
func (s *sqlStorage) FindAll() ([]Experiment, error) {
	rows, err := s.db.Query(
		`SELECT e.name, e.gamma, v.shows, v.rewards
		FROM experiment e INNER JOIN variant v ON v.experiment_name = e.name`,
	)

	if err != nil {
		return nil, err
	}

	var (
		name           string
		gamma          float32
		shows, rewards uint32
		prev           Experiment
		variants       = []Variant{}
		experiments    = []Experiment{}
	)

	for rows.Next() {
		rows.Scan(&name, &gamma, &shows, &rewards)

		if prev != nil && prev.GetName() != name {
			experiments = append(experiments, prev)
			variants = variants[:0]
		}

		variants = append(variants, NewVariant(0, rewards, shows))
		prev = NewExperiment(name, gamma, variants)
	}

	if prev != nil {
		experiments = append(experiments, prev)
	}

	return experiments, nil
}

//
func (s *sqlStorage) Save(experiment Experiment) error {
	if len(experiment.GetVariants()) == 0 {
		return fmt.Errorf("variants should be more or equal to 1")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`REPLACE INTO experiment (name, gamma) VALUES (?, ?)`,
		experiment.GetName(),
		experiment.GetGamma(),
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO variant (id, experiment_name, shows, rewards) VALUES (?, ?, ?, ?)`,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	for i, variant := range experiment.GetVariants() {
		_, err := stmt.Exec(
			i,
			experiment.GetName(),
			variant.GetShows(),
			variant.GetReward(),
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

//
func (s *sqlStorage) Delete(name string) error {
	_, err := s.db.Exec(
		`DELETE FROM experiment WHERE name = ?`,
		name,
	)

	// TODO should i check rowsAffected or it doesn't matter
	return err
}

//
func (s *sqlStorage) IncrementVariant(name string, id, showsDelta, rewardsDelta uint32) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err := increment(tx, name, id, showsDelta, rewardsDelta); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

//
func (s *sqlStorage) Increment(name string, showsDelta, rewardsDelta []uint32) error {
	if len(showsDelta) != len(rewardsDelta) {
		return fmt.Errorf("showsDelta len should be equals to rewardsDelta len")
	}

	if len(showsDelta) == 0 {
		return fmt.Errorf("showsDelta len should be not empty")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for i := 0; i < len(showsDelta); i++ {
		if err := increment(tx, name, uint32(i), showsDelta[i], rewardsDelta[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

//
func increment(tx *sql.Tx, name string, id, showsDelta, rewardsDelta uint32) error {
	if showsDelta == 0 && rewardsDelta == 0 {
		return nil
	}

	rows, err := tx.Exec(
		`UPDATE variant SET shows = shows + ?, rewards = rewards + ? WHERE experiment_name = ? AND id = ?`,
		showsDelta,
		rewardsDelta,
		name,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf(
			"There is no such experiment %s with variant %d",
			name,
			id,
		)
	}

	return nil
}
